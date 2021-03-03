package middleware

import (
	"errors"
	"log"
	"net/http"

	"pokemonb2w/internal/model"
)

const (
	unknownErrorStr = "Unknown error."
)

// Internal method that deals with error messages
func recoverInternal(w http.ResponseWriter) {
	var err error
	r := recover()
	statusCode := http.StatusInternalServerError
	if r != nil {
		log.Println("Recovering from Panic:", r)
		switch t := r.(type) {
		case string:
			err = errors.New(t)
		case model.CustomError:
			err = t
			switch t.ErrorType() {
			case model.ErrorUnprocessableJSON:
				statusCode = http.StatusUnprocessableEntity
			case model.ErrorNotFound:
				statusCode = http.StatusNotFound
			case model.ErrorBadRequest:
				statusCode = http.StatusBadRequest
			case model.ErrorDefault:
				statusCode = http.StatusInternalServerError
			}
		case error:
			err = t
		default:
			err = errors.New(unknownErrorStr)
		}
		log.Printf("Successfuly recovered from panic: %s\n", err.Error())
		http.Error(w, err.Error(), statusCode)
	}
}

// RecoverMiddleware handles panic during requests.
func RecoverMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer recoverInternal(w)
		h.ServeHTTP(w, r)
	})
}
