package api

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/nkansal96/aurora-go/audio"
)

// Based on an input test text, return an audio.File
func GetTTS(text string) (*audio.File, error) {
	// Create GET request
	urlE := url.QueryEscape(fmt.Sprintf("%s%s?text=%s", baseURL, ttsEndpoint, text))
	req, err := http.NewRequest("GET", urlE, nil)
	if err != nil {
		return nil, err
	}

	client := getClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	err = handleError(resp)
	if err != nil {
		return nil, err
	}

	tts := audio.NewFromReader(resp.Body)

	return tts, nil
}
