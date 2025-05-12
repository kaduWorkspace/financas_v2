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
        fmt.Println(err)
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
    for k, _ := range details {
        details[k].Date = months[k]
    }
    self.analizarJurosCompostoService.SetValorFinal(fv)
    self.analizarJurosCompostoService.SetRetornoSobreOInvestimento(post_calcular_cdb.ValorInicial)
    valorizacao := self.analizarJurosCompostoService.GetDiferencaRetorno(self.simularJurosCompostoService.GetValorInicial())
    retorno_sobre_investimento := self.analizarJurosCompostoService.GetRetornoSobreOInvestimento()
    self.simularJurosCompostoService.SetValorGasto()
    self.simularJurosCompostoService.SetValorJurosRendido(self.analizarJurosCompostoService.GetValorFinal())
    valor_gasto := self.simularJurosCompostoService.GetValorGasto()
    dados_tabela := self.analizarJurosCompostoService.AjustarDadosTabela(details, core.EhMobile(ctx.Request().Origin().UserAgent()))
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
func (self FinancasController) Predict(ctx http.Context) http.Response {
    contexto_view := map[string]any{}
    contexto_view["csrf"] = ctx.Request().Session().Get("csrf_token")
    contexto_view["taxa_selic"] = strings.Replace(strconv.FormatFloat(core.GetTaxaSelic(), 'f', 2, 64), ".", ",", -1)
    return ctx.Response().View().Make("predict_wrapper", contexto_view)
}
func (self FinancasController) PredictPost(ctx http.Context) http.Response {
    var post_predict requests.PostPredict
    err := ctx.Request().Bind(&post_predict)
    if err != nil {
        fmt.Println(err.Error())
        return ctx.Response().Header("HX-Redirect", "/?erro=Erro inexperado!").Success().Data("text/plain", []byte(err.Error()))
    }
    ctx.Response().Header("HX-Reswap", "outerHTML")
    errors, err := ctx.Request().ValidateRequest(&post_predict)
    if err != nil {
        fmt.Println(err.Error())
        return ctx.Response().Header("HX-Redirect", "/?erro=Erro inexperado!").Success().Data("text/plain", []byte(err.Error()))
    }
    contexto_view := map[string]any{
        "csrf": ctx.Request().Session().Get("csrf_token"),
        "taxa_selic": strings.Replace(strconv.FormatFloat(core.GetTaxaSelic(), 'f', 2, 64), ".", ",", -1),
    }
    if errors != nil && len(errors.All()) > 0 {
        contexto_view["erros_formulario"] = errors.All()
        fmt.Println(errors)
        ctx.Response().Writer().WriteHeader(http.StatusBadRequest)
        ctx.Response().Header("HX-Retarget", "#form_container")
        return ctx.Response().View().Make("predict_fvs_form", contexto_view)
    }
    self.futureValueOfASeriesService.SetContributionOnFirstDay(post_predict.ContributionOnFirstDay)
    self.futureValueOfASeriesService.SetPeriods(float64(post_predict.Periods))
    self.futureValueOfASeriesService.SetInterestRateDecimal(post_predict.Tax/100)
    var result float64
    if post_predict.InitialValue > 0 {
        result = self.futureValueOfASeriesService.PredictFVWithInitialValue(post_predict.FutureValue, post_predict.InitialValue)
    } else {
        result = self.futureValueOfASeriesService.PredictFV(post_predict.FutureValue)
    }
    contexto_view["valor_final"] = core.FormatarValorMonetario(result)
    contexto_view["aporte_necessario"] = core.FormatarValorMonetario(result)
    contexto_view["valor_inicial"] = core.FormatarValorMonetario(float64(0))
    if post_predict.InitialValue > 0 {
        contexto_view["valor_inicial"] = core.FormatarValorMonetario(post_predict.InitialValue)
    }
    return ctx.Response().View().Make("predict_result", contexto_view)
}
