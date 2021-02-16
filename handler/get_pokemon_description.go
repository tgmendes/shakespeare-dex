package handler

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tgmendes/shakespeare-dex/domain"
	"github.com/tgmendes/shakespeare-dex/web"
)

// PokemonDescriptionResponse represents the response returned from the get pokemon description handler.
type PokemonDescriptionResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (p *shakespeareDexAPI) handleGetPokemonDescription(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	pokemonName := ps.ByName("name")
	if pokemonName == "" {
		web.RespondError(w, "no pokemon name given", http.StatusBadRequest)
		return
	}

	result, err := domain.ShakespeareDescriptionForPokemon(pokemonName, p.pokeapiClient, p.shakespeareClient)
	if err != nil {
		log.Printf("could not get description for %s: %s", pokemonName, err.Error())
		httpErr, ok := err.(*web.HTTPError)
		if ok {
			p.buildClientErrorResponse(w, pokemonName, httpErr)
			return
		}
		web.RespondError(w, "the server encountered a problem and could not complete the request", http.StatusInternalServerError)
		return
	}

	resp := PokemonDescriptionResponse{
		Name:        pokemonName,
		Description: result,
	}

	web.RespondJSON(w, resp, http.StatusOK)
}
