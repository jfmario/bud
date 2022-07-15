package middleware

import (
	"net/http"
	"strings"
)

type methodOverloadMiddleware struct{}

func (m methodOverloadMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		overloadedMethod := r.URL.Query().Get("_method")
		if overloadedMethod != "" {
			r.Method = strings.ToUpper(overloadedMethod)
		}

		next.ServeHTTP(w, r)
	})
}

func MethodOverloadMiddleware() Middleware {
	return methodOverloadMiddleware{}
}
