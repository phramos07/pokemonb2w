package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"pokemonb2w/internal/model"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	maxIdleConnections    = 10
	idleConnectionTimeout = 30
	clientTimeout         = 20

	pokeAPIListPokemonPath = "/pokemon"
	pokeAPIGetPokemonPath  = "/pokemon/{id}"
)

// PokeAPIRequester is responsible for sending requests to the PokeAPI and retrieving data
type PokeAPIRequester interface {
	ListPokemon(offset int, limit int) *model.PokeAPIListPokemonsResponse
	GetPokemon(id int) *model.PokeAPIPokemonResponse

	GetPokemonByURL(url string) *model.PokeAPIPokemonResponse
	GetAreaEncountersByURL(url string) *model.PokeAPILocationAreaEncountersResponse
	GetEvolutionChainsByURL(url string) *model.PokeAPIEvoChainsResponse
	GetPokemonSpeciesByURL(url string) *model.PokeAPIPokemonSpeciesResponse
}

type pokeAPIService struct {
	PokeAPIRequester

	basePath string
	client   http.Client
}

// NewPokeAPIRequester returns new PokeAPIRequester
func NewPokeAPIRequester() PokeAPIRequester {
	return &pokeAPIService{
		basePath: viper.GetString("app.pokeApi.basePath"),
		client: http.Client{
			Transport: &http.Transport{
				MaxIdleConns:       maxIdleConnections,
				IdleConnTimeout:    idleConnectionTimeout * time.Second,
				DisableCompression: true,
			},
			Timeout: clientTimeout * time.Second,
		},
	}
}

func (p *pokeAPIService) ListPokemon(offset int, limit int) *model.PokeAPIListPokemonsResponse {
	var response model.PokeAPIListPokemonsResponse

	// Build request query
	listPokemonURL := fmt.Sprintf("%s%s", p.basePath, pokeAPIListPokemonPath)
	listPokemonURL = strings.TrimRight(listPokemonURL, "/")
	req, err := http.NewRequest(http.MethodGet, listPokemonURL, nil)
	if err != nil {
		log.Panic("Couldn't create request to the PokeAPI.", err.Error())
	}
	q := req.URL.Query()
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	req.URL.RawQuery = q.Encode()

	// Perform request
	p.performRequest(req, &response)

	return &response
}

func (p *pokeAPIService) GetPokemon(id int) *model.PokeAPIPokemonResponse {
	// Build request URL
	getPokemonURL := fmt.Sprintf("%s%s", p.basePath, pokeAPIGetPokemonPath)
	getPokemonURL = strings.Replace(getPokemonURL, "{id}", strconv.Itoa(id), -1)

	return p.GetPokemonByURL(getPokemonURL)
}

func (p *pokeAPIService) GetPokemonByURL(url string) *model.PokeAPIPokemonResponse {
	var response model.PokeAPIPokemonResponse
	url = strings.TrimRight(url, "/")

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Panic("Couldn't create request to the PokeAPI.", err.Error())
	}

	// Perform request
	p.performRequest(req, &response)

	return &response
}

func (p *pokeAPIService) GetAreaEncountersByURL(url string) *model.PokeAPILocationAreaEncountersResponse {
	var response model.PokeAPILocationAreaEncountersResponse
	url = strings.TrimRight(url, "/")

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Panic("Couldn't create request to the PokeAPI.", err.Error())
	}

	// Perform request
	p.performRequest(req, &response.Locations)

	return &response
}

func (p *pokeAPIService) GetEvolutionChainsByURL(url string) *model.PokeAPIEvoChainsResponse {
	var response model.PokeAPIEvoChainsResponse
	url = strings.TrimRight(url, "/")

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Panic("Couldn't create request to the PokeAPI.", err.Error())
	}

	// Perform request
	p.performRequest(req, &response)

	return &response
}

func (p *pokeAPIService) GetPokemonSpeciesByURL(url string) *model.PokeAPIPokemonSpeciesResponse {
	var response model.PokeAPIPokemonSpeciesResponse
	url = strings.TrimRight(url, "/")

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Panic("Couldn't create request to the PokeAPI.", err.Error())
	}

	// Perform request
	p.performRequest(req, &response)

	return &response
}

func (p *pokeAPIService) performRequest(req *http.Request, obj interface{}) {
	// Perform request
	res, err := p.client.Do(req)
	if err != nil {
		log.Panic("Couldn't perform request to the PokeAPI", err.Error())
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusNotFound:
		panic(model.NewrequestError("Resource not found.", model.ErrorNotFound))
	}

	// Read response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Panic("Couldn't read body from request", err.Error())
	}
	err = json.Unmarshal(body, &obj)
	if err != nil {
		log.Panic("Couldn't unmarshal response from request", err.Error())
	}
}
