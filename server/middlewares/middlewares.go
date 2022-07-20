package middlewares

import (
	"backend/services/auth"
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

func IsLoggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := auth.TokenValid(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(context.Background(), "userInfo", claims)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func IsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
		if userInfo["admin"].(bool) {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
		}
	})
}
