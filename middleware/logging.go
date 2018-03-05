package middleware

import (
	"net/http"
	"log"
)

func Logging(logger *log.Logger) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Println("before")
				defer logger.Println("after")
				h.ServeHTTP(w, r)
		})
	}
}
