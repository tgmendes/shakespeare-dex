package shakespeare_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tgmendes/shakespeare-dex/shakespeare"
	"github.com/tgmendes/shakespeare-dex/web"
)

func TestShakespeareAPICall(t *testing.T) {
	var calledPath string
	var contentType string
	var bodyText string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calledPath = r.URL.Path
		b, _ := ioutil.ReadAll(r.Body)
		bodyText = string(b)
		contentType = r.Header.Get("Content-Type")

		w.WriteHeader(http.StatusOK)
		w.Write(mockShakespeareResponse())
	}))
	defer server.Close()

	client := shakespeare.NewClient()
	client.BaseURL = server.URL

	_, err := client.Translate("foo bar")

	assert.NoError(t, err)
	assert.Equal(t, "/translate/shakespeare.json", calledPath)
	assert.Equal(t, "application/x-www-form-urlencoded", contentType)
	assert.Equal(t, "text=foo bar", bodyText)
}

func TestShakespeareTranslation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(mockShakespeareResponse())
	}))
	defer server.Close()

	client := shakespeare.NewClient()
	client.BaseURL = server.URL

	translation, err := client.Translate("foo bar")

	assert.NoError(t, err)
	assert.Equal(t, "Thee did giveth mr. Tim a hearty meal.", translation)
}

func TestAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := shakespeare.NewClient()
	client.BaseURL = server.URL

	_, err := client.Translate("foo bar")

	assert.Error(t, err)
	assert.IsType(t, err, &web.HTTPError{})
}

func TestDecoderError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := shakespeare.NewClient()
	client.BaseURL = server.URL

	_, err := client.Translate("foo bar")

	assert.Error(t, err)
}

func mockShakespeareResponse() []byte {
	resp := `{
		"success": {
			"total": 1
		},
		"contents": {
			"translated": "Thee did giveth mr. Tim a hearty meal.",
			"text": "You gave Mr. Tim a hearty meal.",
			"translation": "shakespeare"
		}
	}`

	return []byte(resp)
}
