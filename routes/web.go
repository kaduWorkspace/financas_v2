package routes

import (
	"goravel/app/http/controllers"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

func Web() {
    financasController := controllers.NewFinancasController()
    facades.Route().Get("/", func(ctx http.Context) http.Response {
        return ctx.Response().View().Make("financas.v2.tmpl", map[string]any{
            "data_inicial": time.Now().Format("2006-01-02"),
            "csrf": ctx.Request().Headers().Get("X-CSRF-Token"),
        })
    })
    facades.Route().Post("/calcular", financasController.CalcularWeb)
}
