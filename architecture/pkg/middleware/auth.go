package middleware

import (
	"arch/ikeppu/github.com/configs"
	"arch/ikeppu/github.com/pkg/jwt"
	"context"
	"net/http"
	"strings"
)

type key string

const (
	ContextEmailKey key = "ContextEmailKey"
)

func Auth(next http.Handler, config *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")

		if !strings.HasPrefix(authorization, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
			return
		}

		token := strings.TrimPrefix(authorization, "Bearer ")

		isValid, data := jwt.NewJWT(config.Auth.Secret).Parse(token)

		if !isValid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
			return
		}

		ctx := context.WithValue(r.Context(), ContextEmailKey, data.Email)

		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}
