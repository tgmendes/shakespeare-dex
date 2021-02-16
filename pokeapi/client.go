package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/tgmendes/shakespeare-dex/domain"

	"github.com/tgmendes/shakespeare-dex/web"
)

const baseURL = "https://pokeapi.co"

// Client is a wrapper around HTTP Client to make calls to the Pokemon API.
type Client struct {
	BaseURL string
	client  *http.Client
}

// NewClient creates a new Pokemon client for the Pokemon API.
func NewClient() *Client {
	c := &http.Client{
		Timeout: 5 * time.Second,
	}

	return &Client{
		BaseURL: baseURL,
		client:  c,
	}
}

// Species represents the species information of a given pokemon.
// It is a stripped down response from /pokeomon-species endpoint -
// any non-relevant fields like ID were skipped.
type Species struct {
	Name              string            `json:"name"`
	FlavorTextEntries []FlavorTextEntry `json:"flavor_text_entries"`
}

// FlavorTextEntry contains the multiple descriptions of a pokemon from different games.
// They can come in different languages.
type FlavorTextEntry struct {
	FlavorText string   `json:"flavor_text"`
	Language   Language `json:"language"`
}

// Language is the ISO name of a language the description is written in.
type Language struct {
	Name string `json:"name"`
}

// GetPokemon retrieves a given pokemon species data and returns it in a domain representation of a Pokemon.
func (c *Client) GetPokemon(name string) (*domain.Pokemon, error) {
	fullPath := fmt.Sprintf("%s%s/%s", c.BaseURL, "/api/v2/pokemon-species", name)

	req, err := http.NewRequest(http.MethodGet, fullPath, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &web.HTTPError{
			StatusCode: resp.StatusCode,
			Err:        fmt.Errorf("error calling %s", fullPath), // TODO - include the error from the API too
		}
	}

	var species Species
	if err = json.NewDecoder(resp.Body).Decode(&species); err != nil {
		return nil, err
	}

	return speciesToPokemon(&species), nil
}

func speciesToPokemon(s *Species) *domain.Pokemon {
	return &domain.Pokemon{
		Name:         s.Name,
		Descriptions: s.getEnglishDescriptions(),
	}
}

func (s *Species) getEnglishDescriptions() []string {
	englishDescriptions := []string{}
	for _, entry := range s.FlavorTextEntries {
		if entry.Language.Name == "en" {
			englishDescriptions = append(englishDescriptions, cleanDescription(entry.FlavorText))
		}
	}

	return englishDescriptions
}

// cleanDescription strips out any newline or \f characters from the given text
func cleanDescription(text string) string {
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\f", "")
	text = strings.ReplaceAll(text, "  ", " ")
	return text
}
