package middleware

import (
	"net/http"

	goravel_http "github.com/goravel/framework/contracts/http"
)
func ValidateCsrfTokenMiddleware() goravel_http.Middleware {
    return func(ctx goravel_http.Context) {
        csrfToken := ctx.Request().Session().Get("csrf_token")
        if csrfToken == nil {
            ctx.Response().Redirect(goravel_http.StatusSeeOther, "/?erro=Erro inexperado")
            ctx.Request().AbortWithStatus(goravel_http.StatusSeeOther)
            return
        }
        csrfTokenEnviado := ctx.Request().Header("csrf_token")
        if csrfTokenEnviado == "" {
            csrfTokenEnviado = ctx.Request().Input("csrf_token")
        }
        if csrfTokenEnviado == "" || csrfTokenEnviado != csrfToken {
            http.Redirect(ctx.Response().Writer(), ctx.Request().Origin(),"/?erro=Erro inexperado!", goravel_http.StatusSeeOther)
            ctx.Request().AbortWithStatus(goravel_http.StatusSeeOther)
            return
        } else {
            ctx.Request().Next()
        }
    }
}
