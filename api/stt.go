package api

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/nkansal96/aurora-go/audio"
)

type STTResponse struct {
	Transcript string `json:"transcript"`
}

func GetSTT(audio *audio.File) (*STTResponse, error) {
	// Pipe the output from the multipart writer so that we don't need to
	// buffer the file in memory before sending it over the network
	r, w := io.Pipe()

	// create a multipart form writing to the pipe
	multi := multipart.NewWriter(w)

	// don't block while copying data
	go func() {
		defer multi.Close()
		part, err := multi.CreateFormFile("audio", "audio.wav")
		if err != nil {
			return
		}
		_, err = part.Write(audio.WAVData())
		if err != nil {
			return
		}
	}()

	// create the request
	req := http.NewRequest("GET", baseURL+sttEndpoint, r)
	if err != nil {
		return nil, err // come back to this
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

	var stt STTResponse
	json.NewDecoder(resp.Body).Decode(&stt)
	return &stt, nil
}
