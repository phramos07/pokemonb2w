package middleware

import (
	"log"
	"net/http"

	"github.com/spf13/viper"
)

const (
	authHeader = "Authorization"
	authEnvVar = "API_KEY"
)

func auth(token string) bool {
	authToken := viper.Get(authEnvVar)

	if authToken != nil {
		return authToken == token
	}
	log.Println("Warning: no api-key set in env vars")

	return true
}

// AuthorizationMiddleware ...
func AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(authHeader)

		if auth(token) {
			log.Printf("Authorized request")
			next.ServeHTTP(w, r)
		} else {
			log.Printf("Unauthorized request.")
			http.Error(w, "Forbidden.", http.StatusForbidden)
		}
	})
}
