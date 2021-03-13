package facade

import (
	"pokemonb2w/internal/model"
	"pokemonb2w/internal/services"
	"sort"
	"sync"
)

// PokemonRequester will be exposed to the controller to work on pokemon requests
type PokemonRequester interface {
	ListPokemon(offset int, limit int, fields []string) *model.PokemonList
	GetPokemon(id int, fields []string) *model.Pokemon
}

type pokemonFacade struct {
	PokemonRequester

	pokeAPIService services.PokeAPIRequester
}

// NewPokemonFacade returns a new PokemonRequester object.
func NewPokemonFacade() PokemonRequester {
	return &pokemonFacade{
		pokeAPIService: services.NewPokeAPIRequester(),
	}
}

// ListPokemon lists all pokemon from the PokeAPI
func (p *pokemonFacade) ListPokemon(offset int, limit int, fields []string) *model.PokemonList {
	pokemonsRes := p.pokeAPIService.ListPokemon(offset, limit)

	var pokemonList *model.PokemonList = &model.PokemonList{
		Results: &model.PokemonListResults{
			Limit:  limit,
			Offset: offset,
			Total:  pokemonsRes.Count,
		},
		Pokemons: make([]*model.Pokemon, 0),
	}

	// Perform all pokemon GET's in parallel
	var wg sync.WaitGroup
	pokemonChannel := make(chan *model.Pokemon, len(pokemonsRes.Results))
	for _, pkm := range pokemonsRes.Results {
		wg.Add(1)
		go func(wg *sync.WaitGroup, url string) {
			defer wg.Done()

			pokemonRes := p.pokeAPIService.GetPokemonByURL(url)
			detailedPkm := p.mountPokemon(pokemonRes, fields)

			pokemonChannel <- detailedPkm
		}(&wg, pkm.URL)

	}

	wg.Wait()
	close(pokemonChannel)

	for p := range pokemonChannel {
		pokemonList.Pokemons = append(pokemonList.Pokemons, p)
	}

	sort.Slice(pokemonList.Pokemons, func(i, j int) bool {
		return pokemonList.Pokemons[i].ID < pokemonList.Pokemons[j].ID
	})

	return pokemonList
}

// GetPokemon retrieves a pokemon by its ID.
func (p *pokemonFacade) GetPokemon(id int, fields []string) *model.Pokemon {
	pokemonPokeAPIResponse := p.pokeAPIService.GetPokemon(id)
	var pokemon *model.Pokemon

	if pokemonPokeAPIResponse != nil {
		pokemon = p.mountPokemon(pokemonPokeAPIResponse, fields)
	}

	return pokemon
}

func (p *pokemonFacade) mountPokemon(pokemonPokeAPIResponse *model.PokeAPIPokemonResponse, fields []string) *model.Pokemon {
	var pokemon *model.Pokemon

	var wg sync.WaitGroup

	empty := (len(fields) == 0)
	if pokemonPokeAPIResponse != nil {
		pokemon = &model.Pokemon{}

		if empty {
			pokemon.ID = pokemonPokeAPIResponse.ID
			pokemon.Name = pokemonPokeAPIResponse.Name
			pokemon.Weight = pokemonPokeAPIResponse.Weight
			pokemon.Height = pokemonPokeAPIResponse.Height

			wg.Add(5)

			go func() {
				defer wg.Done()
				pokemon.Types = p.retrieveTypes(pokemonPokeAPIResponse)
			}()

			go func() {
				defer wg.Done()
				pokemon.LocationAreaEncounters = p.retrieveLocationAreas(pokemonPokeAPIResponse)
			}()

			go func() {
				defer wg.Done()
				pokemon.EvolutionChains = p.retrieveEvoChains(pokemonPokeAPIResponse)
			}()

			go func() {
				defer wg.Done()
				pokemon.Image = p.retrieveImage(pokemonPokeAPIResponse)
			}()

			go func() {
				defer wg.Done()
				pokemon.BaseStats = p.retrieveBaseStats(pokemonPokeAPIResponse)
			}()
		} else {
			for _, field := range fields {
				if field == idFieldName {
					pokemon.ID = pokemonPokeAPIResponse.ID
				} else if field == nameFieldName {
					pokemon.Name = pokemonPokeAPIResponse.Name
				} else if field == weightFieldName {
					pokemon.Weight = pokemonPokeAPIResponse.Weight
				} else if field == heightFieldName {
					pokemon.Height = pokemonPokeAPIResponse.Height
				} else if field == typesFieldName {
					wg.Add(1)
					go func() {
						defer wg.Done()
						pokemon.Types = p.retrieveTypes(pokemonPokeAPIResponse)
					}()
				} else if field == locationAreaEncountersFieldName {
					wg.Add(1)
					go func() {
						defer wg.Done()
						pokemon.LocationAreaEncounters = p.retrieveLocationAreas(pokemonPokeAPIResponse)
					}()
				} else if field == evolutionChainsFieldName {
					wg.Add(1)
					go func() {
						defer wg.Done()
						pokemon.EvolutionChains = p.retrieveEvoChains(pokemonPokeAPIResponse)
					}()
				} else if field == imageFieldName {
					wg.Add(1)
					go func() {
						defer wg.Done()
						pokemon.Image = p.retrieveImage(pokemonPokeAPIResponse)
					}()
				} else if field == baseStatsFieldName {
					wg.Add(1)
					go func() {
						defer wg.Done()
						pokemon.BaseStats = p.retrieveBaseStats(pokemonPokeAPIResponse)
					}()
				}
			}
		}
	}

	wg.Wait()
	return pokemon
}

func (p *pokemonFacade) retrieveTypes(response *model.PokeAPIPokemonResponse) []string {
	var types []string

	for _, v := range response.Types {
		types = append(types, v.Type.Name)
	}

	return types
}

func (p *pokemonFacade) retrieveLocationAreas(response *model.PokeAPIPokemonResponse) []string {
	var locations []string

	locationAreaRes := p.pokeAPIService.GetAreaEncountersByURL(response.LocationAreaEncountersURL)
	for _, v := range locationAreaRes.Locations {
		locations = append(locations, v.LocationArea.Name)
	}

	return locations
}

func (p *pokemonFacade) retrieveEvoChains(response *model.PokeAPIPokemonResponse) [][]string {
	evoChains := make([][]string, 0)

	speciesRes := p.pokeAPIService.GetPokemonSpeciesByURL(response.Species.URL)
	evoChainsRes := p.pokeAPIService.GetEvolutionChainsByURL(speciesRes.EvolutionChain.URL)

	startingChain := make([]string, 0)
	evoChains = *p.fillEvoChainsRec(&evoChains, &startingChain, evoChainsRes.Chain)

	return evoChains
}

func (p *pokemonFacade) fillEvoChainsRec(allPaths *[][]string, current *[]string, chain model.PokeAPIChain) *[][]string {
	if len(chain.EvolvesTo) == 0 {
		*current = append(*current, chain.Species.Name)
		*allPaths = append(*allPaths, *current)

		return allPaths
	}

	*current = append(*current, chain.Species.Name)
	for _, v := range chain.EvolvesTo {
		var newPath []string
		newPath = append(newPath, *current...)
		allPaths = p.fillEvoChainsRec(allPaths, &newPath, v)
	}

	return allPaths
}

func (p *pokemonFacade) retrieveImage(response *model.PokeAPIPokemonResponse) string {
	image := response.Sprites.Versions.GenerationVIII.SwordShield.FrontDefault
	if image != "" {
		return image
	}

	image = response.Sprites.Versions.GenerationVII.UltraSunUltraMoon.FrontDefault
	if image != "" {
		return image
	}

	image = response.Sprites.Versions.GenerationVI.OmegaRubyAlphaSapphire.FrontDefault
	if image != "" {
		return image
	}

	image = response.Sprites.Versions.GenerationVI.XY.FrontDefault
	if image != "" {
		return image
	}

	image = response.Sprites.Versions.GenerationV.BlackWhite.FrontDefault
	if image != "" {
		return image
	}

	image = response.Sprites.Versions.GenerationIV.HeartGoldSoulSilver.FrontDefault
	if image != "" {
		return image
	}

	image = response.Sprites.Versions.GenerationIV.Platinum.FrontDefault
	if image != "" {
		return image
	}

	image = response.Sprites.Versions.GenerationIV.DiamondPearl.FrontDefault
	if image != "" {
		return image
	}

	image = response.Sprites.Versions.GenerationIII.FireRedLeafGreen.FrontDefault
	if image != "" {
		return image
	}

	image = response.Sprites.Versions.GenerationIII.Emerald.FrontDefault
	if image != "" {
		return image
	}

	image = response.Sprites.Versions.GenerationIII.RubySapphire.FrontDefault
	if image != "" {
		return image
	}

	image = response.Sprites.Versions.GenerationII.Crystal.FrontDefault
	if image != "" {
		return image
	}

	image = response.Sprites.Versions.GenerationII.Gold.FrontDefault
	if image != "" {
		return image
	}

	image = response.Sprites.Versions.GenerationII.Silver.FrontDefault
	if image != "" {
		return image
	}

	image = response.Sprites.Versions.GenerationI.Yellow.FrontDefault
	if image != "" {
		return image
	}

	image = response.Sprites.Versions.GenerationI.RedBlue.FrontDefault
	if image != "" {
		return image
	}

	return ""
}

func (p *pokemonFacade) retrieveBaseStats(response *model.PokeAPIPokemonResponse) *model.PokemonBaseStats {
	var baseStats model.PokemonBaseStats

	for _, v := range response.Stats {
		if v.Stat.Name == "hp" {
			baseStats.HP = v.BaseStat
		} else if v.Stat.Name == "attack" {
			baseStats.Attack = v.BaseStat
		} else if v.Stat.Name == "defense" {
			baseStats.Defense = v.BaseStat
		} else if v.Stat.Name == "special-attack" {
			baseStats.SpecialAttack = v.BaseStat
		} else if v.Stat.Name == "special-defense" {
			baseStats.SpecialDefense = v.BaseStat
		} else if v.Stat.Name == "speed" {
			baseStats.Speed = v.BaseStat
		}
	}

	return &baseStats
}
