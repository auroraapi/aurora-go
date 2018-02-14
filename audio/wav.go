package audio

import (
	"io"
	"io/ioutil"
	"encoding/binary"
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

const (
	DefaultNumChannels   = 1
	DefaultSampleRate    = 16000
	DefaultAudioFormat   = 1
	DefaultBitsPerSample = 16
)

// What kind of errors would be possible here?
func NewWAV() (*WAV, error) {
	// create a new default WAV
	b := make([]byte, 0)
	return &WAV{ NumChannels: DefaultNumChannels, SampleRate: DefaultSampleRate, AudioFormat: DefaultAudioFormat, BitsPerSample: DefaultBitsPerSample, audioData: b}, nil
}

// What kind of errors would be possible here?
func NewWAVFromParams(params *WAVParams) (*WAV, error) {
	// create a WAV from the given params
	// use defaults from previous function if any value is 0
	// return error if params are invalid
	if params == nil {
		return NewWAV()
	}
	if params.NumChannels == 0 {
		params.NumChannels = DefaultNumChannels
	}
	if params.SampleRate == 0 {
		params.SampleRate = DefaultSampleRate
	}
	if params.BitsPerSample == 0 {
		params.BitsPerSample = DefaultBitsPerSample
	}
	if binary.Size(params.AudioData) == 0 {
		params.AudioData = make([]byte, 0)
	}
	return &WAV{ NumChannels: params.NumChannels, SampleRate: params.SampleRate, AudioFormat: DefaultAudioFormat, BitsPerSample: params.BitsPerSample, audioData: params.AudioData}, nil

}

// TODO: error checks for this
func NewWAVFromData(data []byte) (*WAV, error) {
	// create a WAV from the given buffer.
	// return error if len(data) < 44
	// extract data from the data according to the spec: http://soundfile.sapp.org/doc/WaveFormat/
	// return error if unexpected data

	// find first data index
	i := 4
	for i < len(data) && data[i-4] != 'd' || data[i-3] != 'a' || data[i-2] != 't' || data[i-1] != 'a' {
		i++
	}

	// ask about this versus the 44 header length
	dataLen := len(data) - i
	if dataLen <= 0 {
		return nil, errors.Error{
			Code:    "One",
			Message: "Received bytes with fewer bytes than the header length",
			Info:    "Critical error",
		}
	}

	numChannels := binary.LittleEndian.Uint16(data[22:23])
	sampleRate := binary.LittleEndian.Uint32(data[24:27])
	bitsPerSample := binary.LittleEndian.Uint16(data[34:35])
	audioData := data[i:]

	return &WAV{ NumChannels: numChannels, SampleRate: sampleRate, AudioFormat: DefaultAudioFormat, BitsPerSample: bitsPerSample, audioData: audioData}, nil
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

func (w *WAV) Data() ([]byte, error) {
	// create header + data (like the function I sent you) based on
	// params stored in w and properties of the data
	// http://soundfile.sapp.org/doc/WaveFormat/
	// remember to set all calculated fields
	// find first data index
	dataLen := binary.Size(w.audioData)
	if dataLen <= 0 {
		return nil, errors.Error{
			Code:    "One",
			Message: "Received bytes with fewer bytes than the header length",
			Info:    "Critical error",
		}
	}

	headerLen := 44
	chunkSize := (dataLen + headerLen - 8)

	// first create the header, then append the rest of the file
	wav := make([]byte, headerLen)
	// RIFF header
	wav[0] = 'R'
	wav[1] = 'I'
	wav[2] = 'F'
	wav[3] = 'F'
	// chunk size
	binary.LittleEndian.PutUint32(wav[4:7], uint32(chunkSize))

	// Format (WAVE)
	wav[8] = 'W'
	wav[9] = 'A'
	wav[10] = 'V'
	wav[11] = 'E'
	// Metadata subchunk ID ("fmt ")
	wav[12] = 'f'
	wav[13] = 'm'
	wav[14] = 't'
	wav[15] = ' '
	// Metadata subchunk size (16)
	wav[16] = 16 
	// Audio format (PCM = 1)
	binary.LittleEndian.PutUint16(wav[20:21], w.AudioFormat) // AUDIO FORMAT
	// Num Channels (Mono = 1)
	binary.LittleEndian.PutUint16(wav[22:23], w.NumChannels) // NUM CHANNELS
	// Sample Rate (16000 Hz)
	binary.LittleEndian.PutUint32(wav[24:27], w.SampleRate) // SAMPLE RATE
	// Byte Rate = SampleRate * NumChannels * BitsPerSample/8 = 32000
	byteRate := uint32(w.SampleRate) * uint32(w.NumChannels) * uint32(w.BitsPerSample) / 8
	binary.LittleEndian.PutUint32(wav[28:31], byteRate)
	// Block Align = NumChannels * BitsPerSample/8
	blockAlign := w.NumChannels * w.BitsPerSample / 8
	binary.LittleEndian.PutUint16(wav[32:33], blockAlign)
	// Bits per sample = 16
	binary.LittleEndian.PutUint16(wav[34:35], w.BitsPerSample)
	// Data subchunk ID ("data")
	wav[36] = 'd'
	wav[37] = 'a'
	wav[38] = 't'
	wav[39] = 'a'
	// Data length
	binary.LittleEndian.PutUint32(wav[40:43], uint32(dataLen))


	return append(wav, w.audioData[0:]...), nil
}
