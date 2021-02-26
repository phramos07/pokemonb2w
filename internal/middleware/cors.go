package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

var corsm *cors.Cors = cors.New(
	cors.Options{
		// Debug: true,

		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},

		AllowCredentials: true,

		AllowedHeaders: []string{
			"Content-Type",
			"api_key",
			"Authorization",
			"Origin",
			"X-Requested-With",
			"Accept",
		},
	})

// CorsMiddleware ...
var CorsMiddleware = corsm.Handler
