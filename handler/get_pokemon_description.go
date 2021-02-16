package handler

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tgmendes/shakespeare-dex/domain"
	"github.com/tgmendes/shakespeare-dex/web"
)

func (p *shakespeareDexAPI) handleGetPokemon(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pokemonName := ps.ByName("name")
	if pokemonName == "" {
		web.RespondError(w, "no pokemon name given", http.StatusBadRequest)
	}

	result, err := domain.ShakespeareDescriptionForPokemon(pokemonName, p.pokeapiClient, p.shakespeareClient)
	if err != nil {
		log.Printf("could not get description for %s: %s", pokemonName, err.Error())
		httpErr, ok := err.(*web.HTTPError)
		if ok {
			p.buildClientErrorResponse(w, pokemonName, httpErr)
		}
		web.RespondError(w, "the server encountered a problem and could not complete the request", http.StatusInternalServerError)
	}

	web.RespondJSON(w, result, http.StatusOK)
}
