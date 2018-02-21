package aurora

import (
	"github.com/nkansal96/aurora-go/api"
)

type Text struct {
	Text string
}

func NewText(text string) *Text {
	return &Text{Text: text}
}

func (t *Text) Speech() (*Speech, error) {
	response, err := api.GetTTS(Config, t.Text)
	return NewSpeech(response), err
}

func (t *Text) Interpret() (*Interpret, error) {
	response, err := api.GetInterpret(Config, t.Text)
	return NewInterpret(response), err
}
