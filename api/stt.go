package api

import (
	"github.com/nkansal96/aurora-go/audio"
)

type STTResponse struct {
	Transcript string `json:"transcript"`
}

func GetSTT(audio *audio.File) (*STTResponse, error) {
	return nil, nil
}