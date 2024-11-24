package controllers

import (
	"goravel/app/core/modules/financas"
	"goravel/app/http/requests"

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
func (_ *FinancasController) Calcular(ctx http.Context) http.Response {
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
    if err := rendimento.Plot("mes"); err != nil {
        println(err.Error())
    }
    return ctx.Response().Success().Json(rendimento)
}
