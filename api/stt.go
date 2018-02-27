package api

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/nkansal96/aurora-go/api/backend"
	"github.com/nkansal96/aurora-go/audio"
	"github.com/nkansal96/aurora-go/config"
)

// STTResponse is the response returned by the API if the speech was
// successfully able to be transcribed.
type STTResponse struct {
	Transcript string `json:"transcript"`
}

// GetSTT queries the API with the provided audio file and returns
// a transcript of the speech.
func GetSTT(c *config.Config, audio *audio.File) (*STTResponse, error) {
	return GetSTTFromStream(c, bytes.NewReader(audio.WAVData()))
}

// GetSTTFromStream queries the API with the provided raw WAV audio stream
// and returns a transcript of the speech.
func GetSTTFromStream(c *config.Config, audio io.Reader) (*STTResponse, error) {
	params := &backend.CallParams{
		Credentials: c.GetCredentials(),
		Method:      "POST",
		Path:        sttEndpoint,
		Body:        audio,
	}

	res, err := c.Backend.Call(params)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	var stt STTResponse
	json.NewDecoder(res.Body).Decode(&stt)
	return &stt, nil
}
