package middleware

import (
	"net/http"

	errs "github.com/Niiaks/campusCart/internal/err"
)

type AuthorizationMiddleware struct{}

func NewAuthorizationMiddleware() *AuthorizationMiddleware {
	return &AuthorizationMiddleware{}
}

// Authorize ensures the request is from an authenticated admin user.
func (am *AuthorizationMiddleware) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetAuthUser(r.Context())
		if user == nil {
			panic(errs.NewUnauthorizedError("not authenticated", false))
		}

		if user.Role != AdminRole {
			panic(errs.NewForbiddenError("admin access required", false))
		}

		next.ServeHTTP(w, r)
	})
}
