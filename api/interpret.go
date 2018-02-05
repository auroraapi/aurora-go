package api

import (
	"encoding/json"
	"net/url"

	"github.com/nkansal96/aurora-go/api/backend"
	"github.com/nkansal96/aurora-go/config"
)

// InterpretResponse is the response returned by the API if the text was
// successfully able to be interpreted
type InterpretResponse struct {
	// Text is the original query
	Text string `json:"text"`

	// Intent is the intent of the user. It is an empty string if the
	// user's intent was unclear
	Intent string `json:"intent"`

	// Entities is a map of the entities in the user's query. The keys
	// are the entity name (like song or location) and the value
	// is the detected value for that entitity
	Entities map[string]string `json:"entities"`
}

// GetInterpret queries the API with the provided text and returns
// the interpreted response
func GetInterpret(c *config.Config, text string) (*InterpretResponse, error) {
	params := &backend.CallParams{
		Credentials: c.GetCredentials(),
		Method:      "GET",
		Path:        interpretEndpoint,
		Query:       url.Values(map[string][]string{"text": []string{text}}),
	}

	res, err := c.Backend.Call(params)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	var i InterpretResponse
	err = json.NewDecoder(res.Body).Decode(&i)
	return &i, err
}
