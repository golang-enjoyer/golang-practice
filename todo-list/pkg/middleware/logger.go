package middleware

import (
	"fmt"
	"net/http"
	"time"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		next.ServeHTTP(w, r)

		endTime := time.Now()
		elapsed := endTime.Sub(startTime)

		fmt.Printf("Method: %s, URL Path: %s - Time Elapsed: %s\n", r.Method, r.URL.Path, elapsed)
	})
}
