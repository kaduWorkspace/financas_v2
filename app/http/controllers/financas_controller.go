package controllers

import (
	"fmt"
	"goravel/app/core"
	"goravel/app/core/modules/financas"
	"goravel/app/http/requests"
	"time"

	"github.com/goravel/framework/contracts/http"
)

type FinancasController struct {
    simularJurosCompostoService financas.SimularJurosComposto
    analizarJurosCompostoService financas.AnalizarResultadoInvestimentoDeJurosComposto
}

func NewFinancasController() *FinancasController {
    simularServiceJC := financas.SimularJurosComposto {}
	return &FinancasController{
        simularJurosCompostoService: simularServiceJC,
        analizarJurosCompostoService: financas.AnalizarResultadoInvestimentoDeJurosComposto {
            JcService : simularServiceJC,
        },
	}
}

func (r *FinancasController) Index(ctx http.Context) http.Response {
    contexto_view := map[string]any{}
    contexto_view["csrf"] = ctx.Request().Session().Get("csrf")
    erro := ctx.Request().Query("erro")
    if  erro != "" {
        contexto_view["panic"] = erro
    }
    return ctx.Response().View().Make("financas.v3.tmpl", contexto_view)
}
func (self *FinancasController) CalcularV2(ctx http.Context) http.Response {
    var post_calcular_cdb requests.PostSimularCdb
    err := ctx.Request().Bind(&post_calcular_cdb)
    if err != nil {
        fmt.Println(err.Error())
        return ctx.Response().Redirect(http.StatusSeeOther, "/?erro=Erro inexperado!")
    }
    errors, err := ctx.Request().ValidateRequest(&post_calcular_cdb)
    contexto_view := map[string]any{}
    if err != nil {
        fmt.Println(err.Error())
        return ctx.Response().Redirect(http.StatusSeeOther, "/?erro=Erro inexperado!")
    }
    if errors != nil && len(errors.All()) > 0 {
        contexto_view["erros_formulario"] = errors.All()
        fmt.Println(errors)
        return ctx.Response().View().Make("financas.v3.tmpl", contexto_view)
    }
    if err := post_calcular_cdb.ValidarData(); err != nil {
        fmt.Println(err.Error())
        return ctx.Response().Redirect(http.StatusSeeOther, "/?erro=Erro inexperado!")
    }
    contexto_view["message"] = "Simulação finalizada!"
    contexto_view["csrf"] = ctx.Request().Session().Get("csrf")
    self.simularJurosCompostoService.SetDatas(post_calcular_cdb.DataInicial, post_calcular_cdb.DataFinal)
    self.simularJurosCompostoService.SetTaxaAnosApartirPeriodoDeDatas()
    self.simularJurosCompostoService.SetDiasDeLiquidesPorAno(post_calcular_cdb.DiasLiquidezPorAno)
    if err := self.simularJurosCompostoService.SetTaxaAnosApartirPeriodoDeDatas(); err != nil {
        contexto_view["panic"] = "Erro inexperado. Tente novamente mais tarde!"
        return ctx.Response().View().Make("financas.v3.tmpl", contexto_view)
    }
    self.simularJurosCompostoService.SetTaxaDeJurosDecimal(post_calcular_cdb.ValorTaxaAnual, financas.PROCENTO_ANUAL)
    var valor_final float64
    var tipo_investimento financas.TIPO_INVESTIMENTO_ENUM
    self.simularJurosCompostoService.SetValorInicial(post_calcular_cdb.ValorInicial)
    if post_calcular_cdb.ValorAporte > 0 {
        tipo_investimento = financas.JC_COM_APORTE_MENSAL_E_VALOR_INICIAL
        self.simularJurosCompostoService.SetValorAporte(post_calcular_cdb.ValorAporte)
        if post_calcular_cdb.ValorInicial > 0 {
            self.simularJurosCompostoService.SetValorInicial(post_calcular_cdb.ValorInicial)
        } else {
            self.simularJurosCompostoService.SetValorInicial(post_calcular_cdb.ValorAporte)
        }
        valor_final = financas.CifAndFvs(self.simularJurosCompostoService.GetValorInicial() ,self.simularJurosCompostoService.GetTaxaDeJurosDecimal(), self.simularJurosCompostoService.GetDiasDeLiquidezPorAno(), self.simularJurosCompostoService.GetTaxaAnos(), self.simularJurosCompostoService.GetValorAporte())
        self.analizarJurosCompostoService.SetValorFinal(valor_final)
        self.analizarJurosCompostoService.SetRetornoSobreOInvestimento(self.simularJurosCompostoService.GetValorInicial())
    }
    if post_calcular_cdb.ValorInicial > 0 && post_calcular_cdb.ValorAporte == 0.0 {
        tipo_investimento = financas.JC_SEM_APORTE
        self.simularJurosCompostoService.SetValorInicial(post_calcular_cdb.ValorInicial)
        valor_final = financas.CompoundInterestFormula(self.simularJurosCompostoService.GetValorInicial(), self.simularJurosCompostoService.GetTaxaDeJurosDecimal(), self.simularJurosCompostoService.GetDiasDeLiquidezPorAno(), self.simularJurosCompostoService.GetTaxaAnos())
        self.analizarJurosCompostoService.SetValorFinal(valor_final)
        self.analizarJurosCompostoService.SetRetornoSobreOInvestimento(self.simularJurosCompostoService.GetValorInicial())
    }

    self.analizarJurosCompostoService.SetTipoInvestimento(tipo_investimento)
    valorizacao := self.analizarJurosCompostoService.GetDiferencaRetorno(self.simularJurosCompostoService.GetValorInicial())
    retorno_sobre_investimento := self.analizarJurosCompostoService.GetRetornoSobreOInvestimento()
    self.simularJurosCompostoService.SetValorGasto()
    self.simularJurosCompostoService.SetValorJurosRendido(valor_final)
    valor_gasto := self.simularJurosCompostoService.GetValorGasto()
    contexto_view["valorizacao"] = core.FormatarValorMonetario(valorizacao)
    contexto_view["valor_investido"] = core.FormatarValorMonetario(valor_gasto)
    contexto_view["valor_inicial"] = core.FormatarValorMonetario(post_calcular_cdb.ValorInicial)
    contexto_view["valor_final"] = core.FormatarValorMonetario(valor_final)
    contexto_view["juros_rendido"] = core.FormatarValorMonetario(self.simularJurosCompostoService.GetValorJurosRendido())
    contexto_view["retorno_sobre_investimento"] = int(retorno_sobre_investimento)
    return ctx.Response().View().Make("financas.v3.tmpl", contexto_view)
}

func (_ *FinancasController) CalcularWeb(ctx http.Context) http.Response {
    var postCalcularJuros requests.PostCalcularJuros
    errors, err := ctx.Request().ValidateRequest(&postCalcularJuros)
    contexto_view := map[string]any{}
    if err != nil {
        fmt.Println(err.Error())
        return ctx.Response().Redirect(http.StatusSeeOther, "/?erro=Erro inexperado!")
    }
    if err := postCalcularJuros.ValidarQuantidades(); err != nil {
        return ctx.Response().Redirect(http.StatusSeeOther, fmt.Sprintf("/?erro=%s", err.Error()))
    }
    if errors != nil && len(errors.All()) > 0 {
        contexto_view["erros_formulario"] = errors.All()
        return ctx.Response().View().Make("financas.v2.tmpl", contexto_view)
    }
    dateLayout := "2006-01-02"
    dataInicial, err := time.Parse(dateLayout, postCalcularJuros.DataInicial)
    if err != nil  {
        fmt.Println(err.Error())
        return ctx.Response().Redirect(http.StatusSeeOther, "/?erro=Erro inexperado!")
    }
    dataFinal, err := time.Parse(dateLayout, postCalcularJuros.DataFinal)
    if err != nil  {
        fmt.Println(err.Error())
        return ctx.Response().Redirect(http.StatusSeeOther, "/?erro=Erro inexperado!")
    }

    postCalcularJuros.DataInicial = dataInicial.Format("02/01/2006")
    postCalcularJuros.DataFinal = dataFinal.Format("02/01/2006")
    if err := postCalcularJuros.ValidarData(); err != nil {
        fmt.Println(err.Error())
        return ctx.Response().Redirect(http.StatusSeeOther, fmt.Sprintf("/?erro=%s", err.Error()))
    }
    service, err := financas.New(postCalcularJuros)
    if err != nil {
        fmt.Println(err.Error())
        return ctx.Response().Redirect(http.StatusSeeOther, "/?erro=Erro inexperado!")
    }
    rendimento := service.Calcular()
    rendimento.SetPorcentagemValorFinalRelativoAValorInicial()
    rendimento.SetMeses()
    rendimento.SetDias()
    rendimento.SetSemestres()
    rendimento.SetDadosProcessados()
    contexto_view["dados_processados"] = rendimento.DadosProcessados
    contexto_view["valorizacao"] = core.FormatarValorMonetario(rendimento.Valorizacao)
    contexto_view["valor_investido"] = core.FormatarValorMonetario(rendimento.Gasto)
    contexto_view["valor_inicial"] = core.FormatarValorMonetario(rendimento.ValorInicial)
    contexto_view["valor_final"] = core.FormatarValorMonetario(rendimento.ValorFinal)
    contexto_view["porcentagem_aumento"] = int(rendimento.PorcentagemValorFinalRelativoAValorInicial)
    fmt.Println(rendimento.PorcentagemValorFinalRelativoAValorInicial, contexto_view["porcentagem_aumento"])
    contexto_view["csrf"] = ctx.Request().Session().Get("csrf")
    json_data, err := rendimento.ToJson()
    if err == nil {
        contexto_view["dados_calculo"] = json_data
    } else {
        fmt.Println(err)
    }
    return ctx.Response().View().Make("financas.v2.tmpl", contexto_view)
}

func (_ *FinancasController) Calcular(ctx http.Context) http.Response {
    return ctx.Response().Json(http.StatusNoContent, http.Json{
        "message": "Route unavailable",
    })
    var postCalcularJuros requests.PostCalcularJuros
    errors, err := ctx.Request().ValidateRequest(&postCalcularJuros)
    if err != nil {
        println(err.Error())
        return ctx.Response().Json(http.StatusInternalServerError, http.Json{
            "message": false,
        })
    }
    if errors != nil && len(errors.All()) > 0 {
        return ctx.Response().Json(400, errors.All())
    }
    if errors := postCalcularJuros.ValidarData(); errors != nil {
        return ctx.Response().Json(400 ,http.Json{
            "Datas": errors.Error(),
        })
    }
    service, err := financas.New(postCalcularJuros)
    if err != nil {
        return ctx.Response().Json(400 ,http.Json{
            "Servico": err.Error(),
        })
    }
    rendimento := service.Calcular()
    rendimento.SetMeses()
    rendimento.SetDias()
    rendimento.SetSemestres()
    return ctx.Response().Success().Json(rendimento)
}
