package routes

import (
	"goravel/app/http/controllers"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

func Web() {
    financasController := controllers.NewFinancasController()
    facades.Route().Fallback(func (ctx http.Context) http.Response {
        return ctx.Response().Redirect(http.StatusSeeOther, "/")
    })
    facades.Route().Get("/", financasController.Index)
    facades.Route().Post("v2/simular-jc", financasController.CalcularV2)
}
