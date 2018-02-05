package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/nkansal96/aurora-go/config"
	"github.com/nkansal96/aurora-go/errors"
)

// URL constants
const (
	baseURL           = "https://api.auroraapi.com"
	sttEndpoint       = "/v1/stt/"
	ttsEndpoint       = "/v1/tts/"
	interpretEndpoint = "/v1/interpret/"
)

// getHeaders creates and returns a map with base headers required for
// all requests
func getHeaders() *http.Header {
	headers := make(http.Header)
	headers.Add("X-Application-Id", config.C.AppID)
	headers.Add("X-Application-Token", config.C.AppToken)
	headers.Add("X-Device-ID", config.C.DeviceID)
	return &headers
}

// getClient creates, configures, and returns a secure https client that
// can be used to make API requests
func getClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
	}
}

// handleError takes an http.Response object and assesses whether or not
// an error occurred. If it did, it returns an error object. Otherwise it
// returns nil
func handleError(r *http.Response) error {
	if r.StatusCode == http.StatusOK {
		return nil
	}

	if r.Header.Get("Content-type") == "application/json" {
		var err errors.APIError
		json.NewDecoder(r.Body).Decode(&err)
		return &err
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return &errors.Error{
		Message: string(body),
		Info:    fmt.Sprintf("An unexpected error occurred with code %d", r.StatusCode),
	}
}
