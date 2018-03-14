package aurora

import (
	"github.com/nkansal96/aurora-go/api"
	"github.com/nkansal96/aurora-go/audio"
	"github.com/nkansal96/aurora-go/errors"
)

const (
	// ListenDefaultLength is the default length of time (in seconds) to listen for.
	ListenDefaultLength = 0.0
	// ListenDefaultSilenceLen is the default amount of silence (in seconds)
	// that the recording framework will allow before stoping.
	ListenDefaultSilenceLen = 0.5
)

// ListenParams configures how the recording framework should listen for
// speech.
type ListenParams struct {
	// Length specifies how long to listen for in seconds. A value of 0
	// means that the recording framework will continue to listen until
	// the specified amount of silence. A value greater than 0 will
	// override any value set to `SilenceLen`
	Length float64
	// SilenceLen is how long of silence (in seconds) will be allowed before
	// automatically stopping. This value is only taken into consideration if
	// `Length` is 0
	SilenceLen float64
}

// NewListenParams creates the default set of ListenParams. You should
// call this function to get the default and then replace the ones you
// want to customize.
func NewListenParams() *ListenParams {
	return &ListenParams{ListenDefaultLength, ListenDefaultSilenceLen}
}

// SpeechHandleFunc is the type of function that is passed to `ContinuouslyListen`.
// It is called every time a speech utterance is decoded and passed to the
// function. If there was an error, that is passed as well. The function must
// return a boolean that indicates whether or not to continue listening (true to
// continue listening, false to stop listening).
type SpeechHandleFunc func(s *Speech, err error) bool

// TextHandleFunc is the type of function that is passed to
// `ContinuouslyListenAndTranscribe`. It is called every time a speech utterance
// is decoded and converted to text. The resulting text object is passed to the
// function. If there was an error, that is passed as well. The function must
// return a boolean that indicates whether or not to continue listening (true to
// continue listening, false to stop listening).
type TextHandleFunc func(t *Text, err error) bool

// Speech represents a user's utterance. It has high-level methods that let
// you operate on the speech (like convert it to text) and also allows you
// to access the underlying audio data to manipulate it, save it, play it, etc.
type Speech struct {
	// Audio is the underlying audio that this struct wraps. You can create
	// a speech object and set this directly if you want to operate on some
	// pre-recorded audio
	Audio *audio.File
}

// NewSpeech creates a Speech object from the given audio file.
func NewSpeech(newAudio *audio.File) *Speech {
	return &Speech{Audio: newAudio}
}

// Text calls the Aurora STT API and converts a user's utterance into
// a text transcription. This is populated into a `Text` object, allowing you
// to chain and combine high-level abstractions.
func (t *Speech) Text() (*Text, error) {
	if t.Audio == nil {
		return nil, errors.NewFromErrorCode(errors.SpeechNilAudio)
	}

	response, err := api.GetSTT(Config, t.Audio)
	if err != nil {
		return nil, err
	}
	return NewText(response.Transcript), nil
}

// `Listen` takes in `ListenParams` and generates a speech object based on those
// parameters by recording from the default input device.
//
// Note that the `ListenParams` is expected to contain values for
// every field, including defaults for fields that you did not want to change.
// To avoid having to do this, you should call `NewListenParams` to obtain an
// instance of `ListenParams` with all of the default filled out, and then over-
// ride them with the ones you want to change. Alternatively, you can pass `nil`
// to simply use the default parameters.
//
// Currently, this function uses the default audio input interface (an option
// to change this will be available at a later time).
func Listen(params *ListenParams) (*Speech, error) {
	if params == nil {
		params = NewListenParams()
	}

	audio, err := audio.NewFileFromRecording(params.Length, params.SilenceLen)
	if err != nil {
		return nil, err
	}
	return &Speech{Audio: audio}, nil
}

// ContinuouslyListen calls `Listen` continuously.

// Note that the `ListenParams` is expected to contain values for every field,
// including defaults for fields that you did not want to change. To avoid having
// to do this, you should call `NewListenParams` to obtain an instance of
// `ListenParams` with all of the default filled out, and then override them with
// the ones you want to change. Alternatively, you can pass `nil` to simply use the
// default parameters.

// This function accepts another function as an argument that is called each time
// a speech utterance is decoded. See the documentation for `SpeechHandleFunc`
// for more information.
func ContinuouslyListen(params *ListenParams, handleFunc SpeechHandleFunc) {
	if params == nil {
		params = NewListenParams()
	}

	for {
		s, err := Listen(params)
		if !handleFunc(s, err) {
			break
		}
	}
}

// ListenAndTranscribe starts listening with the given parameters, except instead
// of waiting for the audio to finish capturing and returning a Speech object,
// it directly streams it to the API, transcribing it in real-time. When the
// transcription completes, this function directly returns a Text object. This
// reduces latency by a significant amount if you already know you want to
// transcribe the audio.
//
// Note that the `ListenParams` is expected to contain values for every field,
// including defaults for fields that you did not want to change. To avoid having
// to do this, you should call `NewListenParams` to obtain an instance of
// `ListenParams` with all of the default filled out, and then override them with
// the ones you want to change. Alternatively, you can pass `nil` to simply use the
// default parameters.
func ListenAndTranscribe(params *ListenParams) (*Text, error) {
	if params == nil {
		params = NewListenParams()
	}

	// create a new recording stream and begin recording. Data will automatically
	// be written to the stream as it becomes available, so we can directly call
	// the API with this stream while audio is recording
	stream := audio.NewRecordingStream(params.Length, params.SilenceLen)
	response, err := api.GetSTTFromStream(Config, stream)
	if err != nil {
		return nil, err
	}
	return NewText(response.Transcript), nil
}

// ContinuouslyListenAndTranscribe is a combination of `ContinuouslyListen`
// and `ListenAndTranscribe`. See the documentation for those two functions to
// understand how it works. The difference is that this handler function receives
// objects of type *Text instead of *Speech. See the documentation for `TextHandleFunc`
// for more information on that.
func ContinuouslyListenAndTranscribe(params *ListenParams, handleFunc TextHandleFunc) {
	if params == nil {
		params = NewListenParams()
	}

	for {
		t, err := ListenAndTranscribe(params)
		if !handleFunc(t, err) {
			break
		}
	}
}
