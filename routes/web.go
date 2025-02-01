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
    facades.Route().Get("v1.0", func(ctx http.Context) http.Response {
        contexto_view := map[string]any{}
        contexto_view["csrf"] = ctx.Request().Session().Get("csrf")
        erro := ctx.Request().Query("erro")
        if  erro != "" {
            contexto_view["panic"] = erro
        }
        contexto_view["csrf"] = ctx.Request().Session().Get("csrf")
        return ctx.Response().View().Make("financas.v2.tmpl", contexto_view)
    })
    facades.Route().Post("v1/simular-jc", financasController.CalcularWeb)
    facades.Route().Post("v2/simular-jc", financasController.CalcularV2)
}
