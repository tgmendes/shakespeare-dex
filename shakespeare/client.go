package shakespeare

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/tgmendes/shakespeare-dex/web"
)

const baseURL = "https://api.funtranslations.com"

// Client is a wrapper around HTTP Client to make calls to the Shakespeare API.
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

// APIResponse is the response from the shakespeare API.
type APIResponse struct {
	Contents Contents `json:"contents"`
}

// Contents contains results of the translation.
type Contents struct {
	Translated string `json:"translated"`
	Text       string `json:"text"`
}

// Translate retrieves a shakespeare translation for a given text.
func (c *Client) Translate(text string) (string, error) {
	fullPath := fmt.Sprintf("%s%s", c.BaseURL, "/translate/shakespeare.json")

	req, err := http.NewRequest(http.MethodPost, fullPath, strings.NewReader(fmt.Sprintf("text=%s", text)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", &web.HTTPError{
			StatusCode: resp.StatusCode,
			Err:        fmt.Errorf("error calling :%s", fullPath), // TODO - include the error from the API too
		}
	}

	var translated APIResponse
	if err = json.NewDecoder(resp.Body).Decode(&translated); err != nil {
		return "", err
	}

	return translated.Contents.Translated, nil
}
