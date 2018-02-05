package api

import (
	"encoding/json"

	"github.com/nkansal96/aurora-go/api/backend"
	"github.com/nkansal96/aurora-go/audio"
	"github.com/nkansal96/aurora-go/config"
)

// STTResponse is the response returned by the API if the speech was
// successfully able to be transcribed
type STTResponse struct {
	Transcript string `json:"transcript"`
}

// GetSTT queries the API with the provided audio file and returns
// a transcript of the speech
func GetSTT(c *config.Config, audio *audio.File) (*STTResponse, error) {
	params := &backend.CallParams{
		Credentials: c.GetCredentials(),
		Method:      "POST",
		Path:        sttEndpoint,
		Files: []backend.MultipartFile{
			backend.MultipartFile{"audio", "audio.wav", audio.WAVData()},
		},
	}

	res, err := c.Backend.CallMultipart(params)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	var stt STTResponse
	json.NewDecoder(res.Body).Decode(&stt)
	return &stt, nil
}
