package controllers

import (
	"fmt"
	"goravel/app/core"
	"goravel/app/core/modules/financas"
	"goravel/app/http/requests"

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
    //fmt.Println("Setando datas: ", self.simularJurosCompostoService.GetDataInicial(), " e ", self.simularJurosCompostoService.GetDataFinal())
    self.simularJurosCompostoService.SetDiasDeLiquidesPorAno(post_calcular_cdb.DiasLiquidezPorAno)
    //fmt.Println("Setando dias liquidez por ano: ", self.simularJurosCompostoService.GetDiasDeLiquidezPorAno())
    if err := self.simularJurosCompostoService.SetTaxaAnosApartirPeriodoDeDatas(); err != nil {
        contexto_view["panic"] = "Erro inexperado. Tente novamente mais tarde!"
        return ctx.Response().View().Make("financas.v3.tmpl", contexto_view)
    }
    //fmt.Println("Setando taxa de anos: ", self.simularJurosCompostoService.GetTaxaAnos())
    self.simularJurosCompostoService.SetTaxaDeJurosDecimal(post_calcular_cdb.ValorTaxaAnual, financas.PROCENTO_ANUAL)
    var valor_final float64
    var tipo_investimento financas.TIPO_INVESTIMENTO_ENUM
    self.simularJurosCompostoService.SetValorInicial(post_calcular_cdb.ValorInicial)
    //fmt.Println("Adicionando valor inicial: ", post_calcular_cdb.ValorInicial)
    if post_calcular_cdb.ValorAporte > 0 {
        //fmt.Println("Usuario definiu um valor para o aporte", post_calcular_cdb.ValorAporte)
        tipo_investimento = financas.JC_COM_APORTE_MENSAL_E_VALOR_INICIAL
        //fmt.Println("Tipo de investimento com valor inicial e aporte mensal")
        self.simularJurosCompostoService.SetValorAporte(post_calcular_cdb.ValorAporte)
        if post_calcular_cdb.ValorInicial > 0 {
            self.simularJurosCompostoService.SetValorInicial(post_calcular_cdb.ValorInicial)
            //fmt.Println("Setando valor inicial como sendo o proprio valor inicial", post_calcular_cdb.ValorInicial)
        } else {
            self.simularJurosCompostoService.SetValorInicial(post_calcular_cdb.ValorAporte)
            //fmt.Println("Setando valor inicial como sendo o valor do aporte", post_calcular_cdb.ValorAporte)
        }
        valor_final = financas.CifAndFvs(self.simularJurosCompostoService.GetValorInicial() ,self.simularJurosCompostoService.GetTaxaDeJurosDecimal(), self.simularJurosCompostoService.GetDiasDeLiquidezPorAno(), self.simularJurosCompostoService.GetTaxaAnos(), self.simularJurosCompostoService.GetValorAporte())
        //fmt.Println("Valor final: ", valor_final)
        self.analizarJurosCompostoService.SetValorFinal(valor_final)
        self.analizarJurosCompostoService.SetRetornoSobreOInvestimento(self.simularJurosCompostoService.GetValorInicial())
        //fmt.Println("Setando retorno sobre o investimento", self.analizarJurosCompostoService.GetRetornoSobreOInvestimento())
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
