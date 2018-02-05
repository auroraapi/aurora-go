package api

import (
	"net/http"
	"net/url"
	"fmt"
	"encoding/json"
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
func GetInterpret(text string) (*InterpretResponse, error) {
	queryStr := url.QueryEscape(text)
	queryUrl := fmt.Sprintf("%s%s?text=%s", baseURL, interpretEndpoint, queryStr)

	req, err := http.NewRequest("GET", queryUrl, nil)
	if err != nil {
		return nil, err // come back to this 
	}

	req.Header = *getHeaders()

	client := getClient()
	res, err := client.Do(req)
	if err != nil {
		return nil, err // come back to this 
	}

	err = handleError(res)
	if err != nil {
		return nil, err
	}

	resStruct := new(InterpretResponse)
	if err := json.NewDecoder(res.Body).Decode(resStruct); err != nil {
		return nil, err
	}

	return resStruct, nil
}
