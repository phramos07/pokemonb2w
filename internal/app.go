// B2W coding challenge API.
//
// OpenAPI doc for the B2W coding challenge.
//
// Terms Of Service:
//
//     Schemes: http, https
//	   BasePath: /v1/api
//     Version: 1.0.0
//     Contact: Supun Muthutantri<fakemail@gmail.com>
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - APIKey:
//
//     SecurityDefinitions:
//     APIKey:
//          type: apiKey
//          name: Authorization
//          in: header
//
// swagger:meta
package main

import (
	"fmt"
	"os"
	"pokemonb2w/internal/control"
	"pokemonb2w/internal/middleware"

	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	// Indirect imports
	_ "github.com/mattn/go-sqlite3"
)

const (
	uiPath          = "./static/"
	swaggerPrefix   = "/swagger/"
	portEnvVariable = "PORT"
	standardPort    = "8080"

	apiPath = "/v1/api"
)

// Base router. Will be used for swaggerUI server
var router *mux.Router = mux.NewRouter()

// API Router. Will be used for the API endpoints and all the middlewares
// of logging, in-memory db and auth.
func getAPIRouter() *mux.Router {
	apiRouter := router.PathPrefix(apiPath).Subrouter()
	return apiRouter
}

// Sets up routes in the API router
func setUpRoutes(apiRouter *mux.Router) {
	control.AddHealthRoutes(apiRouter)
}

// Sets up middlewares in the API router
func setUpRootMiddlewares(root *mux.Router) {
	root.Use(
		middleware.CorsMiddleware,
		middleware.LoggingMiddleware,
		middleware.RecoverMiddleware,
	)
}

func setUpAPIMiddlewares(apiRouter *mux.Router) {
	apiRouter.Use(middleware.AuthorizationMiddleware)
}

// Loads env variables for local development
func loadEnv() {
	// TODO: Use VIPER to load config vars instead of env vars.
	// This STILL NEEDS TO BE CALLED IN DEVELOPMENT. Use env var "env" to
	// decide to load env vars or not.
	log.Println("Loading environment variables for local development")
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found.")
	}
}

// APP's entrypoint
func main() {
	// Load env variables in development env
	loadEnv()

	apiRouter := getAPIRouter()

	// Routes
	setUpRoutes(apiRouter)

	// Middlewares
	setUpRootMiddlewares(router)
	setUpAPIMiddlewares(apiRouter)

	// Serve Swagger UI
	fs := http.FileServer(http.Dir(uiPath))
	router.PathPrefix(swaggerPrefix).Handler(http.StripPrefix(swaggerPrefix, fs))

	// $PORT is defined in the server
	var port string
	port, found := os.LookupEnv(portEnvVariable)

	if !found || port == "" {
		port = standardPort
	}

	log.Printf("Listening on localhost:%s", port)

	log.Fatal(
		http.ListenAndServe(
			fmt.Sprintf(":%s", port),
			router))
}
