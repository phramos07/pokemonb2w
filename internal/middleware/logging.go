package middleware

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

func loggingHandler(w http.ResponseWriter, r *http.Request) {
	headerStr := new(bytes.Buffer)
	fmt.Fprint(headerStr, "{")
	for k, v := range r.Header {
		fmt.Fprintf(headerStr, "%s: %s ", k, v)
	}
	fmt.Fprint(headerStr, "}")
	log.Printf("Request %s %s Header: %s Body: %s", r.Method, r.RequestURI, headerStr.String(), r.Body)
}

// LoggingMiddleware Logs activity in the webapp's requests.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loggingHandler(w, r)
		next.ServeHTTP(w, r)
	})
}
