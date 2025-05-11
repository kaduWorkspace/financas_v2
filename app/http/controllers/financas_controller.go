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
	"github.com/shopspring/decimal"
)

type FinancasController struct {
    simularJurosCompostoService financas.SimularJurosComposto
    analizarJurosCompostoService financas.AnalizarResultadoInvestimentoDeJurosComposto
    compoundInterestService financas.CompoundInterest
    futureValueOfASeriesService financas.FutureValueOfASeries
}

func NewFinancasController() *FinancasController {
    simularServiceJC := financas.SimularJurosComposto {}
    cp := financas.CompoundInterest{}
    fv := financas.FutureValueOfASeries{}
	return &FinancasController{
        compoundInterestService: cp,
        futureValueOfASeriesService: fv,
        simularJurosCompostoService: simularServiceJC,
        analizarJurosCompostoService: financas.AnalizarResultadoInvestimentoDeJurosComposto {
            JcService : simularServiceJC,
        },
	}
}
func (r *FinancasController) Index(ctx http.Context) http.Response {
    return ctx.Response().View().Make("financas")
}

func (r *FinancasController) Simulador(ctx http.Context) http.Response {
    contexto_view := map[string]any{}
    contexto_view["csrf"] = ctx.Request().Session().Get("csrf_token")
    contexto_view["taxa_selic"] = strings.Replace(strconv.FormatFloat(core.GetTaxaSelic(), 'f', 2, 64), ".", ",", -1)
    erro := ctx.Request().Query("erro")
    if  erro != "" {
        contexto_view["panic"] = erro
    }
    return ctx.Response().View().Make("financas_result_wrapper", contexto_view)
}
func (r *FinancasController) Home(ctx http.Context) http.Response {
    return  ctx.Response().View().Make("index")
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
        return ctx.Response().Header("HX-Redirect", "/?erro=Data inválida!").Success().Data("text/plain", []byte(err.Error()))
    }
    contexto_view["message"] = "Simulação finalizada!"
    periods, err := self.compoundInterestService.MonthsBetweenDates(post_calcular_cdb.DataInicial, post_calcular_cdb.DataFinal)
    if err != nil {
        return ctx.Response().Header("HX-Redirect", "/?erro=Data inválida!").Success().Data("text/plain", []byte(err.Error()))
    }
    months, err := self.compoundInterestService.GetDates(post_calcular_cdb.DataInicial, post_calcular_cdb.DataFinal)
    if err != nil {
        return ctx.Response().Header("HX-Redirect", "/?erro=Data inválida!").Success().Data("text/plain", []byte(err.Error()))
    }
    interest, _ := decimal.NewFromFloat(post_calcular_cdb.ValorTaxaAnual).Div(decimal.NewFromInt(100)).Round(16).Float64()
    self.futureValueOfASeriesService.SetInterestRateDecimal(interest)
    self.futureValueOfASeriesService.SetPeriods(float64(periods))
    self.futureValueOfASeriesService.SetContributionAmount(post_calcular_cdb.ValorAporte)
    self.futureValueOfASeriesService.SetContributionOnFirstDay(true)
    fv, details := self.futureValueOfASeriesService.CalculateWithPeriods(post_calcular_cdb.ValorInicial)
    self.analizarJurosCompostoService.SetValorFinal(fv)
    valorizacao := self.analizarJurosCompostoService.GetDiferencaRetorno(self.simularJurosCompostoService.GetValorInicial())
    self.analizarJurosCompostoService.SetRetornoSobreOInvestimento(post_calcular_cdb.ValorInicial)
    retorno_sobre_investimento := self.analizarJurosCompostoService.GetRetornoSobreOInvestimento()
    self.simularJurosCompostoService.SetValorJurosRendido(self.analizarJurosCompostoService.GetValorFinal())
    for k, _ := range details {
        details[k].Date = months[k]
    }
    dados_tabela := self.analizarJurosCompostoService.AjustarDadosTabela(details, core.EhMobile(ctx.Request().Origin().UserAgent()))
    tabela_json, err := json.Marshal(dados_tabela)
    if err == nil {
        contexto_view["tabela_json"] = string(tabela_json)
    } else {
        fmt.Println("ERRO AO GERAR TABELA json: ", err)
    }
    valorInvestidoDecimal := decimal.NewFromInt(int64(periods)).Mul(decimal.NewFromFloat(post_calcular_cdb.ValorAporte)).Add(decimal.NewFromFloat(post_calcular_cdb.ValorInicial)).Round(16)
    valorInvestido, _ := valorInvestidoDecimal.Float64()
    contexto_view["valorizacao"] = core.FormatarValorMonetario(valorizacao)
    contexto_view["valor_investido"] = core.FormatarValorMonetario(valorInvestido)
    contexto_view["valor_inicial"] = core.FormatarValorMonetario(post_calcular_cdb.ValorInicial)
    contexto_view["valor_final"] = core.FormatarValorMonetario(self.analizarJurosCompostoService.GetValorFinal())
    jurosRendido, _ := decimal.NewFromFloat(fv).Sub(valorInvestidoDecimal).Round(16).Float64()
    contexto_view["juros_rendido"] = core.FormatarValorMonetario(jurosRendido)
    contexto_view["retorno_sobre_investimento"] = int(retorno_sobre_investimento)
    contexto_view["taxa_selic"] = strings.Replace(strconv.FormatFloat(self.simularJurosCompostoService.GetTaxaSelic(), 'f', 2, 64), ".", ",", -1)
    contexto_view["tabela"] = dados_tabela
    contexto_view["aporte"] = core.FormatarValorMonetario(post_calcular_cdb.ValorAporte)
    return ctx.Response().View().Make("financas_result", contexto_view)
}
