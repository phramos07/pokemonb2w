package services

import (
	"net/http"
	"pokemonb2w/internal/model"
)

// PokeAPIRequester is responsible for sending requests to the PokeAPI and retrieving data
type PokeAPIRequester interface {
	ListPokemon(offset int) model.PokeAPIListPokemonsResponse
	GetPokemon(id int) model.PokeAPIPokemonResponse

	GetPokemonByURL(url string) model.PokeAPIPokemonResponse
	GetAreaEncountersByURL(url string) model.PokeAPILocationAreaEncountersResponse
	GetEvolutionChainsByURL(url string) model.PokeAPIEvoChainsResponse
	GetPokemonSpeciesByURL(url string) model.PokeAPIPokemonSpeciesResponse
}

type pokeAPIService struct {
	PokeAPIRequester

	defaultClient http.Client
}
