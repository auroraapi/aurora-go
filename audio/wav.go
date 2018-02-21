package audio

import (
	"encoding/binary"
	"io"
	"io/ioutil"
	"math"

	"github.com/nkansal96/aurora-go/errors"
)

// WAV-related constants.
const (
	// DefaultNumChannels is 1 (mono audio)
	DefaultNumChannels = 1
	DefaultSampleRate  = 16000
	// DefaultAudioFormat is 1 (raw, uncompressed PCM waveforms)
	DefaultAudioFormat = 1
	// DefaultBitsPerSample is 16 (2 bytes per sample).
	DefaultBitsPerSample = 16
)

// WAV represents a PCM audio file in the WAV container format. It keeps
// a high-level description of the parameters of the file, along with the
// raw audio bytes, until it needs to be written to a file, stream, or array.
type WAV struct {
	// NumChannels is the number of channels the WAV file has. 1 = mono,
	// 2 = stereo, etc. This affects the block align and also number of
	// bytes per sample: (BitsPerSample / 8) * NumChannels.
	NumChannels uint16
	// SampleRate is the number of samples taken per second.
	SampleRate uint32
	// AudioFormat is the type of Audio that is encoded in the WAV file.
	// In most scenarios, this will be 1 (1 = raw, uncompressed PCM audio),
	// since WAV doesn't support compression.
	AudioFormat uint16
	// BitsPerSample is the width of each sample. 16 bits means each sample
	// is two bytes.
	BitsPerSample uint16

	// audioData is the raw audio data stored in the WAV file.
	audioData []byte
}

// WAVParams are a set of parameters used to create a WAV file. Its fields
// correspond directly to the WAV file.
type WAVParams struct {
	NumChannels   uint16
	SampleRate    uint32
	BitsPerSample uint16
	AudioData     []byte
}

// NewWAV returns a new WAV file from the default parameters. It will never
// return an error.
func NewWAV() *WAV {
	// create a new default WAV
	return &WAV{
		NumChannels:   DefaultNumChannels,
		SampleRate:    DefaultSampleRate,
		AudioFormat:   DefaultAudioFormat,
		BitsPerSample: DefaultBitsPerSample,
		audioData:     make([]byte, 0),
	}
}

// NewWAVFromParams returns a new WAV file from the passed in parameters
// If any of the parameters are 0, then it will be given the default
// values
func NewWAVFromParams(params *WAVParams) *WAV {
	// create a WAV from the given params
	// use defaults from previous function if any value is 0
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
	if params.AudioData == nil || len(params.AudioData) == 0 {
		params.AudioData = make([]byte, 0)
	}
	return &WAV{
		NumChannels:   params.NumChannels,
		SampleRate:    params.SampleRate,
		AudioFormat:   DefaultAudioFormat,
		BitsPerSample: params.BitsPerSample,
		audioData:     params.AudioData,
	}
}

// NewWAVFromData creates a WAV format struct from the given data buffer
// The buffer is broken up into its respective information and that
// information is used to create the WAV format struct
func NewWAVFromData(data []byte) (*WAV, error) {
	// create a WAV from the given buffer.
	// return error if len(data) < 44
	// extract data from the data according to the spec: http://soundfile.sapp.org/doc/WaveFormat/

	// find first data index
	i := 4
	for i < len(data) && data[i-4] != 'd' || data[i-3] != 'a' || data[i-2] != 't' || data[i-1] != 'a' {
		i++
	}

	dataLen := len(data) - i
	if dataLen <= 0 {
		return nil, errors.Error{
			Code:    "One",
			Message: "Received WAV file with empty data",
		}
	}

	// hOff is the header offset. It denotes the start of the header, which is
	// 44 bytes before the data starts.
	hOff := i - 44

	numChannels := binary.LittleEndian.Uint16(data[hOff+22 : hOff+23])
	sampleRate := binary.LittleEndian.Uint32(data[hOff+24 : hOff+27])
	bitsPerSample := binary.LittleEndian.Uint16(data[hOff+34 : hOff+35])
	audioData := data[i:]

	return &WAV{
		NumChannels:   numChannels,
		SampleRate:    sampleRate,
		AudioFormat:   DefaultAudioFormat,
		BitsPerSample: bitsPerSample,
		audioData:     audioData,
	}, nil
}

// NewWAVFromReader takes in a reader and creates a new WAV format
// with the given information.
func NewWAVFromReader(reader io.Reader) (*WAV, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return NewWAVFromData(b)
}

// getRMS is a helper function used to calculate the RMS. This is called
// by TrimSilent which uses RMS to determine whether the sample of audio
// is silent
func getRMS(sampleSize uint16, audioData []byte) float64 {
	sum := 0.0
	for i := uint16(0); i < uint16(len(audioData)); i += sampleSize {
		val := binary.LittleEndian.Uint64(audioData[i:(i + sampleSize)])
		sum += math.Pow(2.0, float64(val))
	}

	return math.Sqrt(sum / (float64(len(audioData) / int(sampleSize))))
}

// TrimSilent is called on a WAV struct to trim the silent portions from the
// ends of the file while leaving a certain amount of padding
func (w *WAV) TrimSilent(threshold float64, padding float64) *WAV {
	// trim silence from the ends of the file, leaving a certain amount of padding
	// numSamples := uint16(len(w.audioData)) * 8 / uint16(w.BitsPerSample)
	sizeOfSample := uint16(w.BitsPerSample / 8)
	rmsSampleCount := uint16((1.0 / 16.0) * float64(w.SampleRate))

	// get max amplitude
	maxPossibleAmp := math.Exp2(float64(w.BitsPerSample))
	maxPossibleVal := maxPossibleAmp / 2.0

	silenceThresh := threshold * maxPossibleVal

	// Trimming the beginning
	N1 := uint16(0)

	for N1 < uint16(len(w.audioData)) {
		// sampleValue := binary.LittleEndian.Uint16(w.audioData[N1:N1+sizeOfSample-1])
		// sumSquares += math.Pow(float64(2), float64(sampleValue))
		// rms := math.Sqrt(sumSquares / float64(numSamples))
		rms := getRMS(sizeOfSample, w.audioData[N1:N1+sizeOfSample*rmsSampleCount])

		N1 += sizeOfSample * rmsSampleCount
		if rms > silenceThresh {
			break
		}
	}

	// Trimming the end
	N2 := uint16(len(w.audioData))
	for N2-uint16(sizeOfSample) >= 0 {
		// sampleValue := binary.LittleEndian.Uint16(w.audioData[N2-sizeOfSample:N2-1])
		// sumSquares += math.Pow(float64(2), float64(sampleValue))
		// rms := math.Sqrt(sumSquares / float64(numSamples))
		rms := getRMS(sizeOfSample, w.audioData[N2-sizeOfSample*rmsSampleCount:N2])

		N2 -= sizeOfSample * rmsSampleCount
		if rms > silenceThresh {
			break
		}
	}
	paddingSamples := uint16(padding*float64(w.SampleRate)) * sizeOfSample
	w.audioData = w.audioData[N1-paddingSamples : N2+paddingSamples-1]
	return w
}

// AddAudioData [... fill in documentation]
func (w *WAV) AddAudioData(d []byte) {
	// add audio data to existing data
	if d != nil && len(d) > 0 {
		w.audioData = append(w.audioData, d...)
	}
}

// AudioData returns the raw audio data
func (w *WAV) AudioData() []byte {
	return w.audioData
}

// Data creates the header and data based on the WAV struct and returns
// a fully formatted WAV file format
func (w *WAV) Data() []byte {
	// find first data index
	dataLen := len(w.audioData)
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
	binary.LittleEndian.PutUint32(wav[16:19], 16)
	// Audio format (PCM = 1)
	binary.LittleEndian.PutUint16(wav[20:21], w.AudioFormat)
	// Num Channels (Mono = 1)
	binary.LittleEndian.PutUint16(wav[22:23], w.NumChannels)
	// Sample Rate (16000 Hz)
	binary.LittleEndian.PutUint32(wav[24:27], w.SampleRate)

	// Byte Rate = SampleRate * NumChannels * BitsPerSample/8 = 32000
	byteRate := w.SampleRate * uint32(w.NumChannels) * uint32(w.BitsPerSample) / 8
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

	return append(wav, w.audioData[0:]...)
}
