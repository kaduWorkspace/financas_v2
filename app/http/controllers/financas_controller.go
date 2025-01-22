package controllers

import (
	"goravel/app/core"
	"goravel/app/core/modules/financas"
	"goravel/app/http/requests"
	"time"

	"github.com/goravel/framework/contracts/http"
)

type FinancasController struct {
	//Dependent services
}

func NewFinancasController() *FinancasController {
	return &FinancasController{
		//Inject services
	}
}

func (r *FinancasController) Index(ctx http.Context) http.Response {
	return nil
}

func (_ *FinancasController) CalcularWeb(ctx http.Context) http.Response {
    var postCalcularJuros requests.PostCalcularJuros
    errors, err := ctx.Request().ValidateRequest(&postCalcularJuros)
    currentDate := time.Now().Format("2006-01-02")
    if err != nil {
        return ctx.Response().View().Make("financas.v2.tmpl", map[string]any{
            "panic": "Erro inexperado!",
            "data_inicial": currentDate,
        })
    }
    if errors != nil && len(errors.All()) > 0 {
        return ctx.Response().View().Make("financas.v2.tmpl", map[string]any{
            "erros_formulario": errors.All(),
            "data_inicial": currentDate,
        })
    }
    dateLayout := "2006-01-02"
    dataInicial, err := time.Parse(dateLayout, postCalcularJuros.DataInicial)
    if err != nil  {
        return ctx.Response().View().Make("financas.v2.tmpl", map[string]any{
            "panic": "Erro inexperado!",
            "data_inicial": currentDate,
        })

    }
    dataFinal, err := time.Parse(dateLayout, postCalcularJuros.DataFinal)
    if err != nil  {
        return ctx.Response().View().Make("financas.v2.tmpl", map[string]any{
            "panic": "Erro inexperado!",
            "data_inicial": currentDate,
        })
    }
    postCalcularJuros.DataInicial = dataInicial.Format("02/01/2006")
    postCalcularJuros.DataFinal = dataFinal.Format("02/01/2006")
    service, err := financas.New(postCalcularJuros)
    if err != nil {
        return ctx.Response().View().Make("financas.v2.tmpl", map[string]any{
            "panic": "Erro inexperado!",
            "data_inicial": currentDate,
        })
    }
    rendimento := service.Calcular()
    rendimento.SetMeses()
    rendimento.SetDias()
    rendimento.SetSemestres()
    rendimentoDadosJson, err := rendimento.ToJson()
    if err != nil {
        return ctx.Response().View().Make("financas.v2.tmpl", map[string]any{
            "panic": "Erro inexperado!",
            "data_inicial": currentDate,
        })
    }
    return ctx.Response().View().Make("financas.v2.tmpl", map[string]any{
        "valorizacao": core.FormatarValorMonetario(rendimento.Valorizacao),
        "valor_investido": core.FormatarValorMonetario(rendimento.Gasto),
        "lucro": core.FormatarValorMonetario(rendimento.Diferenca),
        "valor_inicial": core.FormatarValorMonetario(rendimento.ValorInicial),
        "valor_final": core.FormatarValorMonetario(rendimento.ValorFinal),
        "data_inicial": currentDate,
        "dados_calculo": rendimentoDadosJson,
    })
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
