package domain

import (
	"math/rand"
	"time"
)

// Pokemon represents a single pokemon.
type Pokemon struct {
	// Name is the identifying name of a pokemon.
	Name string

	// Descriptions is a list of different descriptions of the pokemon in English.
	Descriptions []string
}

// PokemonClient is an interface to a client that retrieves pokemon species descriptions.
type PokemonClient interface {
	GetPokemon(name string) (*Pokemon, error)
}

// ShakespeareClient is an interface to a client that retrieves shakespeare translations of a given text.
type ShakespeareClient interface {
	Translate(text string) (string, error)
}

// ShakespeareDescriptionForPokemon will return a description of a given pokemon written in a shakespearean way.
func ShakespeareDescriptionForPokemon(name string, pc PokemonClient, sc ShakespeareClient) (string, error) {
	pokemon, err := pc.GetPokemon(name)
	if err != nil {
		return "", err
	}

	// since our pokemon may have multiple descriptions, we'll just randomly select one
	description := pokemon.getRandomDescription()
	if description == "" {
		return "", nil
	}

	translation, err := sc.Translate(description)
	if err != nil {
		return "", err
	}

	return translation, nil
}

func (p *Pokemon) getRandomDescription() string {
	count := len(p.Descriptions)
	if count == 0 {
		return ""
	}

	rand.Seed(time.Now().UnixNano())

	return p.Descriptions[rand.Intn(count)]
}
