package handler

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tgmendes/shakespeare-dex/domain"
	"github.com/tgmendes/shakespeare-dex/pokeapi"
	"github.com/tgmendes/shakespeare-dex/shakespeare"
	"github.com/tgmendes/shakespeare-dex/web"
)

type shakespeareDexAPI struct {
	pokeapiClient     domain.PokemonClient
	shakespeareClient domain.ShakespeareClient
}

// ShakespeareDexAPI defines and returns the handler for the Pokepeare API.
func ShakespeareDexAPI() http.Handler {
	router := httprouter.New()

	p := shakespeareDexAPI{
		pokeapiClient:     pokeapi.NewClient(),
		shakespeareClient: shakespeare.NewClient(),
	}
	router.GET("/pokemon/:name", p.handleGetPokemon)

	return router
}

// buildErrorResponse will check the HTTP client errors and build an appropriate error response with status code.
//
// For example:
// A 404 from the HTTP client should be returned as 404 from our server, as the given resource couldn't be found.
// A 500 from the HTTP client should be returned as a 502 (bad gateway) from our server.
func (p *shakespeareDexAPI) buildClientErrorResponse(w http.ResponseWriter, resourceName string, err *web.HTTPError) {
	switch err.StatusCode {
	case http.StatusInternalServerError:
		web.RespondError(w, "the server encountered a problem upstream", http.StatusBadGateway)
	case http.StatusBadRequest:
		web.RespondError(w, "the server encountered a problem and could not complete the request", http.StatusInternalServerError)
	case http.StatusTooManyRequests:
		web.RespondError(w, "the server is unavailable due to too many requests", http.StatusServiceUnavailable)
	case http.StatusRequestTimeout:
		web.RespondError(w, fmt.Sprint("the request to upstream server timed out"), http.StatusGatewayTimeout)
	case http.StatusNotFound:
		web.RespondError(w, fmt.Sprintf("the given resource \"%s\" could not be found", resourceName), http.StatusNotFound)
	default:
		web.RespondError(w, fmt.Sprint("the server encountered an unexpected error"), http.StatusInternalServerError)
	}
}
