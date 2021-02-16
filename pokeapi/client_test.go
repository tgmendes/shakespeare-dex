package pokeapi_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tgmendes/shakespeare-dex/pokeapi"
	"github.com/tgmendes/shakespeare-dex/web"
)

func TestPokemonSpeciesAPICall(t *testing.T) {
	var calledPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calledPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		w.Write(mockSpeciesResponse(t))
	}))
	defer server.Close()

	client := pokeapi.NewClient()
	client.BaseURL = server.URL

	_, err := client.GetPokemon("foorizard")

	assert.NoError(t, err)
	assert.Equal(t, "/api/v2/pokemon-species/foorizard", calledPath)
}

func TestGetPokemonSpecies(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(mockSpeciesResponse(t))
	}))
	defer server.Close()

	client := pokeapi.NewClient()
	client.BaseURL = server.URL

	p, err := client.GetPokemon("charizard")

	assert.NoError(t, err)
	assert.NotNil(t, p)

	assert.Equal(t, "charizard", p.Name)
	assert.Len(t, p.Descriptions, 34)
}

func TestAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := pokeapi.NewClient()
	client.BaseURL = server.URL

	_, err := client.GetPokemon("foorizard")

	assert.Error(t, err)
	assert.IsType(t, err, &web.HTTPError{})
}

func TestDecoderError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := pokeapi.NewClient()
	client.BaseURL = server.URL

	_, err := client.GetPokemon("foorizard")

	assert.Error(t, err)
}

func mockSpeciesResponse(t *testing.T) []byte {
	name := "fixtures/pokemon-species.json"
	f, err := ioutil.ReadFile(name)
	if err != nil {
		t.Errorf("failed to read the fixture file %s", name)
		t.FailNow()
	}

	return f
}
