package audio

import (
	"io/ioutil"
	"github.com/nkansal96/aurora-go/errors"
)

type WAV struct {
	NumChannels uint16
	SampleRate uint32
	AudioFormat uint16
	BitsPerSample uint16

	audioData []byte
}

type WAVParams struct {
	NumChannels uint16
	SampleRate uint32
	BitsPerSample uint16
	AudioData []byte
}

func NewWAV() (*WAV, error) {
	// create a new default WAV
	// NumChannels = 1
	// BitsPerSample = 16
	// SampleRate = 16000
}

func NewWAVFromParams(params *WAVParams) (*WAV, error) {
	// create a WAV from the given params
	// use defaults from previous function if any value is 0
	// return error if params are invalid
}

func NewWAVFromData(data []byte) (*WAV, error) {
	// create a WAV from the given buffer.
	// return error if len(data) < 44
	// extract data from the data according to the spec: http://soundfile.sapp.org/doc/WaveFormat/
	// return error if unexpected data
}

func NewWAVFromReader(reader io.Reader) (*WAV, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return NewWAVFromData(b)
}

func (w *WAV) TrimSilent(threshold float64, padding float64) {
	// trim silence from the ends of the file, leaving a certain amount of padding
}

func (w *WAV) AddAudioData(d []byte) {
	// add audio data to existing data
}

func (w *WAV) Data() []byte {
	// create header + data (like the function I sent you) based on
	// params stored in w and properties of the data
	// http://soundfile.sapp.org/doc/WaveFormat/
	// remember to set all calculated fields
}


/////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////
//                                                                 //
//         repurpose the following code into above                 //
//                                                                 //
/////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////

type riffHeader struct {
	ChunkID   [4]byte
	ChunkSize [4]byte
	Format    [4]byte
}

func calculateDataLength(b []byte) ([]byte) {
	i := 4
	for i < len(b) && b[i-4] != 'd' || b[i-3] != 'a' || b[i-2] != 't' || b[i-1] != 'a' {
		i++
	}

	dataLen := len(b) - i
	return dataLen
}

func newRiffHeader(chunkSize uint32) *riffHeader{
	riffHeaderLen := 12
	wav := make([]byte, riffHeaderLen)

	// RIFF header
	wav[0] = 'R'
	wav[1] = 'I'
	wav[2] = 'F'
	wav[3] = 'F'
	// chunk size
	wav[4] = byte((chunkSize >> 0) & 0xFF)
	wav[5] = byte((chunkSize >> 8) & 0xFF)
	wav[6] = byte((chunkSize >> 16) & 0xFF)
	wav[7] = byte((chunkSize >> 24) & 0xFF)
	// Format (WAVE)
	wav[8] = 'W'
	wav[9] = 'A'
	wav[10] = 'V'
	wav[11] = 'E'

	return &riffHeader{ ChunkID: wav[0:3], ChunkSize: wav[4:7], Format: wav[8:11] }
}

func MakeValidWAV(b []byte) ([]byte, *errors.Error) *WAV {
	dataLen := calculateDataLength(b)

	if dataLen <= 0 {
		return nil, errors.FromErrorCodeWithInfo(errors.TTSResponseInvalidAudioData, fmt.Sprintf("Received bytes with length %d and header length of %d", len(b), i))
	}

	headerLen := 44
	chunkSize := (dataLen + headerLen - 8)

	riffHeader := newRiffHeader(chunkSize)
	
	//Currently only creates the header
	return WAV{ RiffHeader: riffHeader }
}
