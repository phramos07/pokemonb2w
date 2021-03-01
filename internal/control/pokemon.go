package control

import (
	"encoding/json"
	"log"
	"net/http"
	"pokemonb2w/internal/facade"
	"pokemonb2w/internal/model"
	"strconv"

	"github.com/gorilla/mux"
)

// PokemonResponse response model
//
// This is used for returning a response with a pokemon inside the body
//
// swagger:response pokemonResponse
type PokemonResponse struct {
	// in: body
	Pokemon *model.Pokemon `json:"pokemon"`
}

// GetPokemonParams params for retrieving pokemon by ID.
//
// This is used for retrieving pokemon by ID.
//
// swagger:parameters getPokemon
type GetPokemonParams struct {
	// in: path
	// required: true
	ID int `json:"id"`
}

// swagger:operation GET /pokemon/{id} pokemon getPokemon
// ---
// summary: Retrieves a pokemon by id.
// description: Retrieves details about a pokemon by its id.
// responses:
//   '200':
//     "$ref": "#/responses/pokemonResponse"
//   '204':
//     description: No pokemon found.
//     schema:
//       type: string
func getPokemon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pokeID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Panic(model.NewrequestError("Bad ID param type.", model.ErrorBadRequest))
	}

	pokemon := facade.GetPokemon(pokeID)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(pokemon)
	if err != nil {
		panic(err)
	}
}

// AddPokemonRoutes Adds all routes from pokemon controller to the router.
func AddPokemonRoutes(r *mux.Router) {
	r.HandleFunc(pokemonByIDPath, getPokemon).Methods(http.MethodGet, http.MethodOptions)
}
