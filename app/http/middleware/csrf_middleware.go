package middleware

import (
	"fmt"
    net_http "net/http"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/gorilla/csrf"
)
func CsrfMiddleware() http.Middleware {
    app_key := fmt.Sprintf("%s",facades.Config().Env("APP_KEY"))
    csrgGorillaMiddleware := csrf.Protect([]byte(app_key), csrf.Secure(true))
	return func(ctx http.Context) {
        handler := csrgGorillaMiddleware(net_http.HandlerFunc(func(w net_http.ResponseWriter, r *net_http.Request) {
			token := csrf.Token(r)
            if(token == "") {
                fmt.Println("Sem token")
                ctx.Response().Redirect(303, "/")
            } else {
                fmt.Println("Com token")
                ctx.Request().Session().Put("csrf", token)
                ctx.Request().Session().Save()
                ctx.Request().Next()
            }
        }))
        handler.ServeHTTP(ctx.Response().Writer(), ctx.Request().Origin())
	}
}
