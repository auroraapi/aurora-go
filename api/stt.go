package api

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/nkansal96/aurora-go/audio"
)

// STTResponse is the response returned by the API if the speech was
// successfully able to be transcribed
type STTResponse struct {
	Transcript string `json:"transcript"`
}

// GetSTT uses the API to convert an audio file to speech
func GetSTT(audio *audio.File) (*STTResponse, error) {
	// Pipe the output from the multipart writer so that we don't need to
	// buffer the file in memory before sending it over the network
	r, w := io.Pipe()

	// create a multipart form writing to the pipe
	multi := multipart.NewWriter(w)

	// don't block while copying data
	go func() {
		defer w.Close()
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
	req, err := http.NewRequest("GET", baseURL+sttEndpoint, r)
	if err != nil {
		return nil, err // come back to this
	}

	client := getClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err // come back to this
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
