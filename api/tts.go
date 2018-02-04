package api

import (
	"net/http"

	"github.com/nkansal96/aurora-go/audio"
)

func GetTTS(string text) (*audio.File, error) {
	// Create GET request
	req, err := http.NewRequest("GET", baseURL+ttsEndpoint, &text)
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

	// Take data from server, put into audio.File
	var tts audio.File
}
