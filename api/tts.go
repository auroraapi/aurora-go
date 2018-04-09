package api

import (
	"net/url"

	"github.com/auroraapi/aurora-go/api/backend"
	"github.com/auroraapi/aurora-go/audio"
	"github.com/auroraapi/aurora-go/config"
)

// GetTTS calls the TTS API given some text and returns an *audio.File
// with the audio from converting the text to speech.
func GetTTS(c *config.Config, text string) (*audio.File, error) {
	params := &backend.CallParams{
		Credentials: c.GetCredentials(),
		Method:      "GET",
		Path:        ttsEndpoint,
		Query:       url.Values(map[string][]string{"text": []string{text}}),
	}

	res, err := c.Backend.Call(params)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	return audio.NewFileFromReader(res.Body)
}
