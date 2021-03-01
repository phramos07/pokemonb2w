package facade

import (
	"pokemonb2w/internal/model"
	"pokemonb2w/internal/services"
)

var pokeAPIService services.PokeAPIRequester

// StartPokeAPIService starts singleton for pokeAPI service
func StartPokeAPIService() {
	if pokeAPIService == nil {
		pokeAPIService = services.NewPokeAPIRequester()
	}
}

// GetPokemon ...
func GetPokemon(id int) *model.Pokemon {
	pokemonPokeAPIResponse := pokeAPIService.GetPokemon(id)

	pokemon := &model.Pokemon{
		ID:                     pokemonPokeAPIResponse.ID,
		Name:                   pokemonPokeAPIResponse.Name,
		Weight:                 pokemonPokeAPIResponse.Weight,
		Height:                 pokemonPokeAPIResponse.Height,
		Types:                  retrieveTypes(pokemonPokeAPIResponse),
		LocationAreaEncounters: retrieveLocationAreas(pokemonPokeAPIResponse),
		EvolutionChains:        retrieveEvoChains(pokemonPokeAPIResponse),
		Image:                  retrieveImage(pokemonPokeAPIResponse),
		BaseStats:              retrieveBaseStats(pokemonPokeAPIResponse),
	}

	return pokemon
}

func retrieveTypes(response *model.PokeAPIPokemonResponse) []string {
	var types []string

	for _, v := range response.Types {
		types = append(types, v.Type.Name)
	}

	return types
}

func retrieveLocationAreas(response *model.PokeAPIPokemonResponse) []string {
	var locations []string

	locationAreaRes := pokeAPIService.GetAreaEncountersByURL(response.LocationAreaEncountersURL)
	for _, v := range locationAreaRes.Locations {
		locations = append(locations, v.LocationArea.Name)
	}

	return locations
}

func retrieveEvoChains(response *model.PokeAPIPokemonResponse) [][]string {
	evoChains := make([][]string, 0)

	speciesRes := pokeAPIService.GetPokemonSpeciesByURL(response.Species.URL)
	evoChainsRes := pokeAPIService.GetEvolutionChainsByURL(speciesRes.EvolutionChain.URL)

	startingChain := make([]string, 0)
	evoChains = *fillEvoChainsRec(&evoChains, &startingChain, evoChainsRes.Chain)

	return evoChains
}

func fillEvoChainsRec(allPaths *[][]string, current *[]string, chain model.PokeAPIChain) *[][]string {
	if len(chain.EvolvesTo) == 0 {
		*current = append(*current, chain.Species.Name)
		*allPaths = append(*allPaths, *current)

		return allPaths
	}

	*current = append(*current, chain.Species.Name)
	for _, v := range chain.EvolvesTo {
		var newPath []string
		newPath = append(newPath, *current...)
		allPaths = fillEvoChainsRec(allPaths, &newPath, v)
	}

	return allPaths
}

func retrieveImage(response *model.PokeAPIPokemonResponse) string {
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

	return "no_image_retrieved"
}

func retrieveBaseStats(response *model.PokeAPIPokemonResponse) *model.PokemonBaseStats {
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
