package routes

import (
	"github.com/goravel/framework/facades"

	"goravel/app/http/controllers"
)

func Api() {
	financasController := controllers.NewFinancasController()
    facades.Route().Post("/calcular-juros", financasController.Calcular)
}
