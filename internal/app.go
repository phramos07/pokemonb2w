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
	"pokemonb2w/internal/control"
	"pokemonb2w/internal/facade"
	"pokemonb2w/internal/middleware"

	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

const (
	standardPort = "8080"

	viperConfigFileName = "config"
	viperConfigFileType = "yml"
	viperConfigPath     = "./internal/"
)

// Base router. Will be used for swaggerUI server
var router *mux.Router = mux.NewRouter()

// API Router. Will be used for the API endpoints and all the middlewares
// of logging, in-memory db and auth.
func getAPIRouter() *mux.Router {
	apiRouter := router.PathPrefix(
		viper.GetString("app.api.path"),
	).Subrouter()

	return apiRouter
}

// Sets up routes in the API router
func setUpRoutes(apiRouter *mux.Router) {
	control.AddAppInfoRoute(apiRouter)
	control.AddHealthRoutes(apiRouter)
	control.AddPokemonRoutes(apiRouter)
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
func loadConfig() {
	viper.SetConfigName(viperConfigFileName)
	viper.SetConfigType(viperConfigFileType)
	viper.AddConfigPath(viperConfigPath)

	err := viper.ReadInConfig()

	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		log.Println("Config file not found.")
	} else if err != nil {
		log.Fatal(
			"Couldn't read config file:",
			err.Error(),
		)
	}

	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()
}

func serveSwagger() {
	uiPath := viper.GetString("app.swagger.uiPath")
	fs := http.FileServer(http.Dir(uiPath))

	swaggerPrefix := viper.GetString("app.swagger.prefix")
	router.PathPrefix(swaggerPrefix).Handler(
		http.StripPrefix(swaggerPrefix, fs),
	)
}

func startServices() {
	facade.StartPokeAPIService()
}

// APP's entrypoint
func main() {
	// Load config variables
	loadConfig()

	apiRouter := getAPIRouter()

	// Routes
	setUpRoutes(apiRouter)

	// Middlewares
	setUpRootMiddlewares(router)
	setUpAPIMiddlewares(apiRouter)

	// Serve Swagger UI
	serveSwagger()

	// Start services
	startServices()

	// $PORT is defined in the server
	port := viper.Get("port")
	if port == nil || port == "" {
		port = standardPort
	}

	log.Printf("Listening on localhost:%s", port)
	log.Fatal(
		http.ListenAndServe(
			fmt.Sprintf(":%s", port),
			router))
}
