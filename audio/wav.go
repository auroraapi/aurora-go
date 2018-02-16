package audio

import (
	"io"
	"math"
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
	if len(params.AudioData) == 0 {
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

	dataLen := len(data) - i
	if dataLen <= 0 {
		return nil, errors.Error{
			Code:    "One",
			Message: "Received bytes with fewer bytes than the header length",
			Info:    "Critical error",
		}
	}

	headerOffset := i - 44 

	numChannels := binary.LittleEndian.Uint16(data[headerOffset+22:headerOffset+23])
	if numChannels > 65535 {
		return nil, errors.Error{
			Code:    "One",
			Message: "The number of channels must be less than or equal to 65535 but greater than 0",
			Info:    "Critical error",
		}
	}
	sampleRate := binary.LittleEndian.Uint32(data[headerOffset+24:headerOffset+27])
	// Error check for sample rate?
	bitsPerSample := binary.LittleEndian.Uint16(data[headerOffset+34:headerOffset+35])
	if bitsPerSample % 8 != 0 {
		return nil, errors.Error{
			Code:    "One",
			Message: "The bits per sample must be divisible by 8",
			Info:    "Critical error",
		}
	}
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

func getRMS(rmsSampleCount uint16, begin uint16, end uint16, audioData []byte) float64{
	sampleValue := binary.LittleEndian.Uint16(audioData[begin:end])
	sumSquares := math.Pow(float64(2), float64(sampleValue))
	rms := math.Sqrt(sumSquares / float64(rmsSampleCount))
	return rms
}

func (w *WAV) TrimSilent(threshold float64, padding float64) *WAV {
	// trim silence from the ends of the file, leaving a certain amount of padding
	// numSamples := uint16(len(w.audioData)) * 8 / uint16(w.BitsPerSample)
	sizeOfSample := uint16(w.BitsPerSample / 8)
	rmsSampleCount := uint16(1/16 * w.SampleRate)

	// get max amplitude 
	max_possible_amp := uint16(math.Exp2(float64(w.BitsPerSample)))
	max_possible_val := max_possible_amp / 2 

	var silence_thresh float64 = float64(threshold) * float64(max_possible_val)

	// Trimming the beginning
	N1 := uint16(0)

	for N1 < uint16(len(w.audioData)) {
		// sampleValue := binary.LittleEndian.Uint16(w.audioData[N1:N1+sizeOfSample-1])
		// sumSquares += math.Pow(float64(2), float64(sampleValue))
		// rms := math.Sqrt(sumSquares / float64(numSamples))
		rms := getRMS(rmsSampleCount, N1, N1 + sizeOfSample * rmsSampleCount - 1, w.audioData)

		N1 += sizeOfSample * rmsSampleCount
		if rms > silence_thresh {
			break
		}
	}
	
	// Trimming the end
	N2 := uint16(len(w.audioData))
	for N2 - uint16(sizeOfSample) >= 0 {
		// sampleValue := binary.LittleEndian.Uint16(w.audioData[N2-sizeOfSample:N2-1])
		// sumSquares += math.Pow(float64(2), float64(sampleValue))
		// rms := math.Sqrt(sumSquares / float64(numSamples))
		rms := getRMS(rmsSampleCount, N2 - sizeOfSample * rmsSampleCount, N2 - 1, w.audioData)

		N2 -= sizeOfSample * rmsSampleCount
		if rms > silence_thresh {
			break
		}
	}
	padding_samples := uint16((uint32(padding) * w.SampleRate)) * sizeOfSample
	w.audioData = w.audioData[N1-padding_samples:N2+padding_samples-1]
	return w
}


func (w *WAV) AddAudioData(d []byte) (*WAV, error) {
	// add audio data to existing data
	if d == nil {
		return nil, errors.Error {
			Code:    "One",
			Message: "The received audio data was nil",
			Info:    "Critical error",
		}
	}
	w.audioData = append(w.audioData, d...)
	return w, nil
}

func (w *WAV) Data() ([]byte, error) {
	// create header + data (like the function I sent you) based on
	// params stored in w and properties of the data
	// http://soundfile.sapp.org/doc/WaveFormat/
	// remember to set all calculated fields
	// find first data index
	dataLen := len(w.audioData)
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
	binary.LittleEndian.PutUint32(wav[16:19], 16)
	// Audio format (PCM = 1)
	binary.LittleEndian.PutUint16(wav[20:21], w.AudioFormat) // AUDIO FORMAT
	// Num Channels (Mono = 1)
	binary.LittleEndian.PutUint16(wav[22:23], w.NumChannels) // NUM CHANNELS
	// Sample Rate (16000 Hz)
	binary.LittleEndian.PutUint32(wav[24:27], w.SampleRate) // SAMPLE RATE
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

	return append(wav, w.audioData[0:]...), nil
}
