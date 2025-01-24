package http

import (
	"goravel/app/http/middleware"

	"github.com/goravel/framework/contracts/http"
    session_mid "github.com/goravel/framework/session/middleware"
)

type Kernel struct {
}

// The application's global HTTP middleware stack.
// These middleware are run during every request to your application.
func (kernel Kernel) Middleware() []http.Middleware {
	return []http.Middleware{
        session_mid.StartSession(),
        middleware.CsrfMiddleware(),
    }
}
