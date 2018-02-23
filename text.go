package aurora

import (
	"github.com/nkansal96/aurora-go/api"
)

// Text encapsulates some text, whether it is obtained from STT, a user input,
// or generated programmatically, and allows high-level operations to be
// conducted and chained on it (like converting to speech, or calling Interpret)
type Text struct {
	// Text is the actual text that this object encapsulates
	Text string
}

// NewText creates a Text object from the given text
func NewText(text string) *Text {
	return &Text{Text: text}
}

// Speech calls the Aurora STT service on the text encapulated in this object
// and converts it to a `Speech` object. Further operations can then be done
// on it, such as saving to file or speaking the resulting audio.
func (t *Text) Speech() (*Speech, error) {
	response, err := api.GetTTS(Config, t.Text)
	if err != nil {
		return nil, err
	}
	return NewSpeech(response), nil
}

// Interpret calls the Aurora Interpret service on the text encapsulated in this
// object and converts it to an `Interpret` object, which contains the results
// from the APIcalls
func (t *Text) Interpret() (*Interpret, error) {
	response, err := api.GetInterpret(Config, t.Text)
	if err != nil {
		return nil, err
	}
	return NewInterpret(response), nil
}
