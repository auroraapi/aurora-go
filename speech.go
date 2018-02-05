package aurora

import (
	"github.com/nkansal96/aurora-go/api"
	"github.com/nkansal96/aurora-go/audio"
	"github.com/nkansal96/aurora-go/errors"
)

type Speech struct {
	Audio *audio.File
}

const (
	ListenDefaultLength     = 0.0
	ListenDefaultSilenceLen = 1.0
)

type ListenParams struct {
	Length     float64
	SilenceLen float64
}

func NewSpeech(newAudio *audio.File) *Speech {
	return &Speech{Audio: newAudio}
}

func NewListenParams() *ListenParams {
	return &ListenParams{ListenDefaultLength, ListenDefaultSilenceLen}
}

func (t *Speech) Text() (*Text, error) {
	if t.Audio == nil {
		return nil, errors.Error{
			Code:    "One",
			Message: "Audio file from Speech object is empty",
			Info:    "Critical error",
		}
	}

	response, err := api.GetSTT(Config, t.Audio)

	if err != nil {
		return nil, errors.Error{
			Code:    "One",
			Message: "Text from speech is empty",
			Info:    "Critical error",
		}
	}

	return NewText(response.Transcript), err
}

func ContinouslyListen(params *ListenParams) (chan *Speech, chan bool) {
	if params == nil {
		params = NewListenParams()
	}

	listenChannel := make(chan (*Speech))
	doneChannel := make(chan bool)

	go func() {
		defer close(listenChannel)
		for {
			select {
			case listenChannel <- listen(params):
			case <-doneChannel:
				return
			}
		}
	}()

	return listenChannel, doneChannel
}

func listen(params *ListenParams) *Speech {
	if params == nil {
		params = NewListenParams()
	}

	return &Speech{Audio: audio.NewFromRecording(params.Length, params.SilenceLen)}
}
