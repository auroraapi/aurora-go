package audio

import (
	"encoding/binary"
	"io"
	"io/ioutil"
	"os"

	_ "github.com/nkansal96/aurora-go/errors"

	"github.com/gordonklaus/portaudio"
)

const (
	BufSize      = 1 << 10
	MaxThresh    = 1 << 14
	SilentThresh = 1 << 10
	SampleRate   = 16000
	NumChannels  = 1
)

// File is an audio file
type File struct {
	AudioData *WAV
}

// Writes the audio data to a file
func (f *File) WriteToFile(filename string) error {
	return ioutil.WriteFile(filename, f.AudioData.Data(), 0644)
}

// Pad adds silence to both the beginning and end of the audio data. Silence
// is specified in seconds.
func (f *File) Pad(seconds float64) {
	// calculate number of bytes needed to pad given amount of seconds
	bytes := float64(f.AudioData.BitsPerSample/8) * float64(f.AudioData.SampleRate) * float64(f.AudioData.NumChannels) * seconds
	padding := make([]byte, int(bytes))

	// copy the WAV parameters, set initial data to the left padding
	newWav := NewWAVFromParams(&WAVParams{
		NumChannels:   f.AudioData.NumChannels,
		SampleRate:    f.AudioData.SampleRate,
		BitsPerSample: f.AudioData.BitsPerSample,
		AudioData:     padding,
	})

	// add the original data and ther right padding
	newWav.AddAudioData(f.AudioData.AudioData())
	newWav.AddAudioData(padding)

	// set the audio data to the new wav
	f.AudioData = newWav
}

// PadLeft adds silence to the beginning of the audio data
func (f *File) PadLeft(seconds float64) {
	// calculate number of bytes needed to pad given amount of seconds
	bytes := float64(f.AudioData.BitsPerSample/8) * float64(f.AudioData.SampleRate) * float64(f.AudioData.NumChannels) * seconds
	padding := make([]byte, int(bytes))

	// copy the WAV parameters, set initial data to the left padding
	newWav := NewWAVFromParams(&WAVParams{
		NumChannels:   f.AudioData.NumChannels,
		SampleRate:    f.AudioData.SampleRate,
		BitsPerSample: f.AudioData.BitsPerSample,
		AudioData:     padding,
	})

	// add the original data
	newWav.AddAudioData(f.AudioData.AudioData())

	// set the audio data to the new wav
	f.AudioData = newWav
}

//PadRight adds silence to the end of the audio data
func (f *File) PadRight(seconds float64) {
	// calculate number of bytes needed to pad given amount of seconds
	bytes := float64(f.AudioData.BitsPerSample/8) * float64(f.AudioData.SampleRate) * float64(f.AudioData.NumChannels) * seconds
	padding := make([]byte, int(bytes))

	// add the padding to the right
	f.AudioData.AddAudioData(padding)
}

//Trims all silence from the audio data
func (F *File) TrimSilence() {
	// TODO: calibrate constants
	f.AudioData.TrimSilent(0.10, 0.25)
}

// Play plays the Audio File to the default output
func (f *File) Play() error {
	// initialize the underlying APIs for audio transmission
	portaudio.Initialize()
	defer portaudio.Terminate()

	// create a buffer for audio to be put into
	// hardcode 16-bit sample until i figure out a better way to do this
	bufLen := int(BufSize * f.AudioData.NumChannels)
	buf := make([]uint16, bufLen)

	// create the audio stream to write to
	stream, err := portaudio.OpenDefaultStream(0, int(f.AudioData.NumChannels), float64(f.AudioData.SampleRate), BufSize, buf)
	if err != nil {
		return err
	}

	defer stream.Close()
	defer stream.Stop()
	stream.Start()

	//Making the assumption that no audio file will be under 64 bytes long
	data := f.AudioData.AudioData()
	for i := 0; i <= len(data); i += bufLen * 2 {
		// need to convert each 2-bytes in [i, i+buffLen*2] to little endian uint16
		for j := 0; j < bufLen; j++ {
			buf[j] = binary.LittleEndian.Uint16(data[i+(j*2) : i+(j*2+1)])
		}
		err := stream.Write()
		if err != nil {
			// should we do something here? we ignore it in python
		}
	}

	return nil
}

// NewFromRecording creates a new File by recording from the default input stream.
// length specifies the maximum length of the recording in seconds. silenceLen
// specifies how long in seconds to automatically stop after when consecutive
// silence is detected.
func NewFromRecording(length float64, silenceLen float64) (*File, error) {
	portaudio.Initialize()
	defer portaudio.Terminate()

	buf := make([]uint16, BufSize)
	data := make([]uint16, 0)

	stream, err := portaudio.OpenDefaultStream(NumChannels, 0, SampleRate, BufSize, buf)
	if err != nil {
		return nil, err
	}

	defer stream.Close()
	defer stream.Stop()
	stream.Start() // check err

	silentFor := 0.0
	for {
		err := stream.Read()
		if err != nil {
			// should we do something here? we ignore it in python
		}

		data = append(data, buf...)

		if IsSilent(buf) {
			silentFor += float64(len(buf)) / SampleRate
		}

		if length == 0 && silentFor > silenceLen {
			break
		}

		if length > 0 && len(data) > int(length*SampleRate) {
			break
		}
	}

	audioData := make([]byte, 2*len(data))
	for i := 0; i < len(data); i += 2 {
		audioData[i] = byte(data[i] & 0xFF)
		audioData[i+1] = byte((data[i] >> 8) & 0xFF)
	}

	return &File{
		AudioData: NewWAVFromParams(&WAVParams{
			NumChannels:   NumChannels,
			SampleRate:    SampleRate,
			BitsPerSample: 16,
			AudioData:     audioData,
		}),
	}, nil
}

//NewFromBytes creates a new Audio File from WAV data
func NewFromBytes(b []byte) (*File, error) {
	wav, err := NewWAVFromData(b)
	return &File{wav}, err
}

//Creates a new Audio File from an io.Reader
func NewFromReader(r io.Reader) (*File, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return NewFromBytes(data)
}

//Creates a new Audio File from an os.File
func NewFromFile(f *os.File) (*File, error) {
	return NewFromReader(f)
}

//Creates a new Audio File from a given filename
func NewFromFileName(f string) (*File, error) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return NewFromBytes(data)
}

//Returns the wav data contained in the audio file
func (f *File) WAVData() []byte {
	return f.AudioData.Data()
}

//Determines whether an audio slice is silent or not
func IsSilent(audio []uint16) bool {
	var max uint16 = audio[0]
	for _, value := range audio {
		if value > max {
			max = value
		}
	}
	return max < SilentThresh
}
