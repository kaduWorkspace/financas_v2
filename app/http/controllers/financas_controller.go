package controllers

import (
	"encoding/json"
	"fmt"
	"goravel/app/core"
	"goravel/app/core/modules/financas"
	"goravel/app/http/requests"
	"strconv"
	"strings"

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
    contexto_view["csrf"] = ctx.Request().Session().Get("csrf_token")
    contexto_view["taxa_selic"] = strings.Replace(strconv.FormatFloat(core.GetTaxaSelic(), 'f', 2, 64), ".", ",", -1)
    erro := ctx.Request().Query("erro")
    if  erro != "" {
        contexto_view["panic"] = erro
    }
    return ctx.Response().View().Make("financas", contexto_view)
}
func (self *FinancasController) CalcularV2(ctx http.Context) http.Response {
    var post_calcular_cdb requests.PostSimularCdb

    err := ctx.Request().Bind(&post_calcular_cdb)
    if err != nil {
        fmt.Println(err.Error())
        return ctx.Response().Redirect(http.StatusSeeOther, "/?erro=Erro inexperado!")
    }
    ctx.Response().Header("HX-Reswap", "outerHTML")
    errors, err := ctx.Request().ValidateRequest(&post_calcular_cdb)
    contexto_view := map[string]any{
        "csrf": ctx.Request().Session().Get("csrf_token"),
    }
    if err != nil {
        return ctx.Response().Header("HX-Redirect", "/?erro=Erro inexperado!").Success().Data("text/plain", []byte(err.Error()))
    }
    if errors != nil && len(errors.All()) > 0 {
        contexto_view["erros_formulario"] = errors.All()
        fmt.Println(errors)
        ctx.Response().Writer().WriteHeader(http.StatusBadRequest)
        ctx.Response().Header("HX-Retarget", "#form_container")
        return ctx.Response().View().Make("financas_form", contexto_view)
    }
    if err := post_calcular_cdb.ValidarData(); err != nil {
        return ctx.Response().Header("HX-Redirect", "/?erro=Erro inexperado!").Success().Data("text/plain", []byte(err.Error()))
    }
    contexto_view["message"] = "Simulação finalizada!"
    self.simularJurosCompostoService.SetDatas(post_calcular_cdb.DataInicial, post_calcular_cdb.DataFinal)
    self.simularJurosCompostoService.SetDiasDeLiquidesPorAno(post_calcular_cdb.DiasLiquidezPorAno)
    if err := self.simularJurosCompostoService.SetTaxaAnosApartirPeriodoDeDatas(); err != nil {
        return ctx.Response().Header("HX-Redirect", "/?erro=Erro inexperado!").Success().Data("text/plain", []byte(err.Error()))
    }
    if self.simularJurosCompostoService.GetTaxaAnos() < 1 {
        // escolheu 6 meses
        self.simularJurosCompostoService.SetTaxaAnos(0.5)
    } else {
        self.simularJurosCompostoService.SetTaxaAnos(float64(int(self.simularJurosCompostoService.GetTaxaAnos())))
    }
    self.simularJurosCompostoService.SetTaxaDeJurosDecimal(post_calcular_cdb.ValorTaxaAnual, financas.PORCENTO_ANUAL)
    var tipo_investimento financas.TIPO_INVESTIMENTO_ENUM
    self.simularJurosCompostoService.SetValorInicial(post_calcular_cdb.ValorInicial)
    if post_calcular_cdb.ValorAporte > 0 {
        tipo_investimento = financas.JC_COM_APORTE_MENSAL_E_VALOR_INICIAL
        self.simularJurosCompostoService.SetValorAporte(post_calcular_cdb.ValorAporte)

        // Define o valor inicial: se ValorInicial > 0, usa ele; caso contrário, usa o ValorAporte
        valorInicial := post_calcular_cdb.ValorInicial
        if valorInicial <= 0 {
            valorInicial = post_calcular_cdb.ValorAporte
        }
        self.simularJurosCompostoService.SetValorInicial(valorInicial)
    } else if post_calcular_cdb.ValorInicial > 0 {
        tipo_investimento = financas.JC_SEM_APORTE
        self.simularJurosCompostoService.SetValorInicial(post_calcular_cdb.ValorInicial)
    }
    self.simularJurosCompostoService.SetTaxaMesesRangeData()
    resultado, err := financas.FutureValueOfASeriesMonthly(self.simularJurosCompostoService.GetValorInicial(), self.simularJurosCompostoService.GetTaxaDeJurosDecimal(), self.simularJurosCompostoService.GetDiasDeLiquidezPorAno(),self.simularJurosCompostoService.GetValorAporte(), float64(int(self.simularJurosCompostoService.GetTaxaMeses())), true, self.simularJurosCompostoService.GetDataInicial())
    if err != nil {
        fmt.Println(err)
        return ctx.Response().Header("HX-Redirect", "/?erro=Erro inexperado!").Success().Data("text/plain", []byte(err.Error()))
    }
    self.analizarJurosCompostoService.SetValorFinal(resultado[len(resultado)-1].Acumulado)
    self.analizarJurosCompostoService.SetRetornoSobreOInvestimento(self.simularJurosCompostoService.GetValorInicial())
    self.analizarJurosCompostoService.SetTipoInvestimento(tipo_investimento)
    valorizacao := self.analizarJurosCompostoService.GetDiferencaRetorno(self.simularJurosCompostoService.GetValorInicial())
    retorno_sobre_investimento := self.analizarJurosCompostoService.GetRetornoSobreOInvestimento()
    self.simularJurosCompostoService.SetValorGasto()
    self.simularJurosCompostoService.SetValorJurosRendido(self.analizarJurosCompostoService.GetValorFinal())
    valor_gasto := self.simularJurosCompostoService.GetValorGasto()
    dados_tabela := self.analizarJurosCompostoService.AjustarDadosTabela(resultado, core.EhMobile(ctx.Request().Origin().UserAgent()))
    tabela_json, err := json.Marshal(dados_tabela)
    if err == nil {
        contexto_view["tabela_json"] = string(tabela_json)
    } else {
        fmt.Println("ERRO AO GERAR TABELA json: ", err)
    }
    contexto_view["valorizacao"] = core.FormatarValorMonetario(valorizacao)
    contexto_view["valor_investido"] = core.FormatarValorMonetario(valor_gasto)
    contexto_view["valor_inicial"] = core.FormatarValorMonetario(post_calcular_cdb.ValorInicial)
    contexto_view["valor_final"] = core.FormatarValorMonetario(self.analizarJurosCompostoService.GetValorFinal())
    contexto_view["juros_rendido"] = core.FormatarValorMonetario(self.simularJurosCompostoService.GetValorJurosRendido())
    contexto_view["retorno_sobre_investimento"] = int(retorno_sobre_investimento)
    contexto_view["taxa_selic"] = strings.Replace(strconv.FormatFloat(self.simularJurosCompostoService.GetTaxaSelic(), 'f', 2, 64), ".", ",", -1)
    contexto_view["tabela"] = dados_tabela
    contexto_view["aporte"] = core.FormatarValorMonetario(post_calcular_cdb.ValorAporte)
    return ctx.Response().View().Make("financas_result", contexto_view)
}
