package middleware

import (
	"fmt"
	"goravel/app/core"

	"github.com/goravel/framework/contracts/http"
)

func CreateCsrfTokenMiddleware() http.Middleware {
	return func(ctx http.Context) {
        if ctx.Request().Session().Has("csrf_token") {
            ctx.Request().Next()
            return
        } else {
            novo_token := core.GerarTokenApartirDeAppKey()
            fmt.Println("Token adicionado a sessao")
            ctx.Request().Session().Put("csrf_token", novo_token)
            ctx.Request().Next()
        }
	}
}
