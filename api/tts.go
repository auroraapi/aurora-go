package api

import (
	"io"
	"net/http"

	"github.com/nkansal96/aurora-go/audio"
)

// Based on an input test text, return an audio.File
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
	var audioReader io.Reader
	audioReader.Read(resp.Body)
	tts = audio.NewFromReader(audioReader)

	return &tts, nil
}
