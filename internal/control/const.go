package control

const (
	healthPath  = "/health"
	apiInfoPath = "/app-info"

	pokemonByIDPath = "/pokemon/{id}"
	pokemonPath     = "/pokemon"

	defaultListOffset = 0
	defaultListLimit  = 20

	offsetQueryParam = "offset"
	limitQueryParam  = "limit"
	fieldsQueryParam = "fields"

	fieldsSeparator = ","
)
