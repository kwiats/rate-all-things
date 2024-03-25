package middleware

import (
	"net/http"
	"tit/internal/auth"
)

func AuthMiddleware(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")
		isValid, _ := auth.CheckValidationToken(token)
		if !isValid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		handlerFunc(w, r)

	}
}
