package model

// Pokemon it's a model that describes a Pokemon.
type Pokemon struct {
	ID                     int               `json:"id,omitempty"`
	Name                   string            `json:"name,omitempty"`
	Image                  string            `json:"image,omitempty"`
	Types                  []string          `json:"types,omitempty"`
	LocationAreaEncounters []string          `json:"locationAreaEncounters,omitempty"`
	EvolutionChains        [][]string        `json:"evolutionChains,omitempty"`
	Weight                 int               `json:"weight,omitempty"`
	Height                 int               `json:"height,omitempty"`
	BaseStats              *PokemonBaseStats `json:"baseStats,omitempty"`
}

// PokemonBaseStats is a model that describes a Pokemon's base stats.
type PokemonBaseStats struct {
	Speed          int `json:"speed"`
	HP             int `json:"hp"`
	Attack         int `json:"attack"`
	Defense        int `json:"defense"`
	SpecialAttack  int `json:"special-attack"`
	SpecialDefense int `json:"special-defense"`
}

// PokemonList is a model that describes a list of Pokemon obtained from the list pokemon request
type PokemonList struct {
	Pokemons []*Pokemon          `json:"pokemons"`
	Results  *PokemonListResults `json:"_results"`
}

// PokemonListResults is a model that describes metadata from list pokemon request
type PokemonListResults struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

// PokeAPIPokemonResponse describes the response obtained from PokeAPI for a Pokemon
type PokeAPIPokemonResponse struct {
	ID     int    `json:"ID"`
	Name   string `json:"name"`
	Height int    `json:"height"`
	Weight int    `json:"weight"`

	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`

	LocationAreaEncountersURL string `json:"location_area_encounters"` // URL to retrive location area

	Species struct {
		URL string `json:"url"` // URL to retrieve pokemon species
	} `json:"species"`

	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		}
	} `json:"stats"`

	Sprites struct {
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					FrontDefault string `json:"front_default"`
				} `json:"red-blue"`
				Yellow struct {
					FrontDefault string `json:"front_default"`
				} `json:"yellow"`
			} `json:"generation-i"`

			GenerationII struct {
				Crystal struct {
					FrontDefault string `json:"front_default"`
				} `json:"crystal"`
				Gold struct {
					FrontDefault string `json:"front_default"`
				} `json:"gold"`
				Silver struct {
					FrontDefault string `json:"front_default"`
				} `json:"silver"`
			} `json:"generation-ii"`

			GenerationIII struct {
				Emerald struct {
					FrontDefault string `json:"front_default"`
				} `json:"emerald"`
				RubySapphire struct {
					FrontDefault string `json:"front_default"`
				} `json:"ruby-sapphire"`
				FireRedLeafGreen struct {
					FrontDefault string `json:"front_default"`
				} `json:"firered-leafgreen"`
			} `json:"generation-iii"`

			GenerationIV struct {
				DiamondPearl struct {
					FrontDefault string `json:"front_default"`
				} `json:"diamond-pearl"`
				Platinum struct {
					FrontDefault string `json:"front_default"`
				} `json:"platinum"`
				HeartGoldSoulSilver struct {
					FrontDefault string `json:"front_default"`
				} `json:"heartgold-soulsilver"`
			} `json:"generation-iv"`

			GenerationV struct {
				BlackWhite struct {
					FrontDefault string `json:"front_default"`
				} `json:"black-white"`
			} `json:"generation-v"`

			GenerationVI struct {
				XY struct {
					FrontDefault string `json:"front_default"`
				} `json:"x-y"`
				OmegaRubyAlphaSapphire struct {
					FrontDefault string `json:"front_default"`
				} `json:"omegaruby-alphasapphire"`
			} `json:"generation-vi"`

			GenerationVII struct {
				UltraSunUltraMoon struct {
					FrontDefault string `json:"front_default"`
				} `json:"ultra-sun-ultra-moon"`
			} `json:"generation-vii"`

			GenerationVIII struct {
				SwordShield struct {
					FrontDefault string `json:"front_default"`
				} `json:"sword-shield"`
			} `json:"generation-viii"`
		} `json:"versions"`
	} `json:"sprites"`
}

// PokeAPIChain describes a recursive evolution chain in a PokeAPI evolution-chain response
type PokeAPIChain struct {
	EvolvesTo []PokeAPIChain `json:"evolves_to"`

	Species struct {
		Name string `json:"name"`
	} `json:"species"`
}

// PokeAPIEvoChainsResponse describes a evolution chain received by PokeAPI
type PokeAPIEvoChainsResponse struct {
	Chain PokeAPIChain `json:"chain"`
}

// PokeAPILocationAreaEncountersResponse describes a response from Location Area Encounters request
type PokeAPILocationAreaEncountersResponse struct {
	Locations []struct {
		LocationArea struct {
			Name string `json:"name"`
		} `json:"location_area"`
	}
}

// PokeAPIListPokemonsResponse describes a response from PokeAPI pokemon list request
type PokeAPIListPokemonsResponse struct {
	Count int `json:"count"`

	Results []struct {
		Name string `json:"name"`
		URL  string `json:"url"` // URL to retrieve details from that Pokemon
	} `json:"results"`
}

// PokeAPIPokemonSpeciesResponse describes a response from PokeAPI pokemon species request
type PokeAPIPokemonSpeciesResponse struct {
	EvolutionChain struct {
		URL string `json:"url"` // URL to retrieve evolution chain
	} `json:"evolution_chain"`
}
