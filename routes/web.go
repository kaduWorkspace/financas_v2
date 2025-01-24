package routes

import (
	"goravel/app/http/controllers"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

func Web() {
    financasController := controllers.NewFinancasController()
    facades.Route().Fallback(func (ctx http.Context) http.Response {
        return ctx.Response().Redirect(http.StatusSeeOther, "/")
    })
    facades.Route().Get("/", func(ctx http.Context) http.Response {
        contexto_view := map[string]any{}
        contexto_view["data_inicial"] = time.Now().Format("2006-01-02")
        contexto_view["csrf"] = ctx.Request().Session().Get("csrf")
        erro := ctx.Request().Query("erro")
        if  erro != "" {
            contexto_view["panic"] = erro
        }
        return ctx.Response().View().Make("financas.v2.tmpl", contexto_view)
    })
    facades.Route().Post("/calcular", financasController.CalcularWeb)
}
