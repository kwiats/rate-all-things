package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var bodyBytes []byte
		if r.Body != nil {
			bodyBytes, _ = io.ReadAll(r.Body)
		}

		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		log.Printf("[%s] %s \n%s", r.Method, r.RequestURI, string(bodyBytes))
		next.ServeHTTP(w, r)
	})
}
