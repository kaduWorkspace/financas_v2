package middleware

import (
	"net/http"

	goravel_http "github.com/goravel/framework/contracts/http"
)
type CommonPostRequest struct {
    csrf_token string `form:"csrf_token"`
}
func ValidateCsrfTokenMiddleware() goravel_http.Middleware {
    return func(ctx goravel_http.Context) {
        csrfToken := ctx.Request().Session().Get("csrf_token")
        if csrfToken == nil {
            ctx.Response().Redirect(goravel_http.StatusSeeOther, "/?erro=Erro de segurança!")
            ctx.Request().AbortWithStatus(goravel_http.StatusSeeOther)
            return
        }
        csrfTokenEnviado := ctx.Request().Header("csrf_token")
        if csrfTokenEnviado == "" {
            csrfTokenEnviado = ctx.Request().Input("csrf_token")
        }
        if csrfTokenEnviado == "" || csrfTokenEnviado != csrfToken {
            http.Redirect(ctx.Response().Writer(), ctx.Request().Origin(),"/?erro=Erro de segurança!", goravel_http.StatusSeeOther)
            ctx.Request().AbortWithStatus(goravel_http.StatusSeeOther)
            return
        } else {
            ctx.Request().Next()
        }
    }
}
