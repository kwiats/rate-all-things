package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var logBody string

		if r.Header.Get("Content-Type") != "" && strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			logBody = "<binary_file>"
		} else {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusInternalServerError)
				log.Printf("Error reading request body: %v\n", err)
				return
			}

			logBody = string(bodyBytes)
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		log.Printf("[%s] %s \n%s", r.Method, r.RequestURI, logBody)
		next.ServeHTTP(w, r)
	})
}
