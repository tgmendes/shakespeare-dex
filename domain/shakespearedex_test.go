package domain_test

import (
	"errors"
	"testing"

	"github.com/tgmendes/shakespeare-dex/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetDescription(t *testing.T) {
	pokemonClient := new(mockPokemonClient)
	pokemonClient.On("GetPokemon", mock.Anything).Return(&domain.Pokemon{
		Name: "foobar",
		Descriptions: []string{
			"something",
			"something else",
		},
	}, nil)

	shakespeareClient := new(mockShakespeareClient)
	shakespeareClient.On("Translate", mock.Anything).Return("thee something", nil)

	description, err := domain.ShakespeareDescriptionForPokemon("foobar", pokemonClient, shakespeareClient)

	assert.NoError(t, err)
	assert.Equal(t, "thee something", description)
}

func TestGetEmptyDescriptions(t *testing.T) {
	pokemonClient := new(mockPokemonClient)
	pokemonClient.On("GetPokemon", mock.Anything).Return(&domain.Pokemon{
		Name:         "foobar",
		Descriptions: []string{},
	}, nil)

	shakespeareClient := new(mockShakespeareClient)
	shakespeareClient.On("Translate", mock.Anything).Return("thee something", nil)

	description, err := domain.ShakespeareDescriptionForPokemon("foobar", pokemonClient, shakespeareClient)

	assert.NoError(t, err)
	assert.Equal(t, "", description)
	// if no descriptions are available, we don't want to call the shakespeare server
	shakespeareClient.AssertNotCalled(t, "Translate")
}

func TestPokemonClientError(t *testing.T) {
	pokemonClient := new(mockPokemonClient)
	pokemonClient.On("GetPokemon", mock.Anything).Return(nil, errors.New("some client error"))

	_, err := domain.ShakespeareDescriptionForPokemon("foobar", pokemonClient, new(mockShakespeareClient))

	assert.Error(t, err)
}

func TestShakespeareClientError(t *testing.T) {
	pokemonClient := new(mockPokemonClient)
	pokemonClient.On("GetPokemon", mock.Anything).Return(&domain.Pokemon{
		Name: "foobar",
		Descriptions: []string{
			"something",
			"something else",
		},
	}, nil)

	shakespeareClient := new(mockShakespeareClient)
	shakespeareClient.On("Translate", mock.Anything).Return(nil, errors.New("some client error"))

	_, err := domain.ShakespeareDescriptionForPokemon("foobar", pokemonClient, shakespeareClient)

	assert.Error(t, err)
}

// mocks for testing purposes
type mockPokemonClient struct {
	mock.Mock
}

func (m *mockPokemonClient) GetPokemon(name string) (*domain.Pokemon, error) {
	args := m.Called(name)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Pokemon), args.Error(1)
}

type mockShakespeareClient struct {
	mock.Mock
}

func (m *mockShakespeareClient) Translate(text string) (string, error) {
	args := m.Called(text)

	if args.Get(0) == nil {
		return "", args.Error(1)
	}
	return args.Get(0).(string), args.Error(1)
}
