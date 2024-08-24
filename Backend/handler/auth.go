package handler

import (
	"context"
	"net/http"
	"strings"
)

// HTTPInterceptor 拦截器
func HTTPInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "No authorization token provided", http.StatusUnauthorized)
				return
			}

			// 解析 Bearer token
			token := strings.TrimPrefix(authHeader, "Bearer ")

			username, err := ValidateTokenAndGetUsername(token)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// 将 username 添加到请求的上下文中
			ctx := context.WithValue(r.Context(), "username", username)
			r = r.WithContext(ctx)
			h(w, r)

		})

}
