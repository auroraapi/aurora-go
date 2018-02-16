package audio

import  (
	"os"
	"io"
	"math"
	"github.com/gordonklaus/portaudio"
	"azul3d.org/engine/audio"
)

const (
	BufSize uint16 = Math.pow(2,10)
	MaxThresh uint16 = Math.pow(2,14)
	SilentThresh uint16 = Math.pow(2,10)
	NumChannels uint16 = 1
	Rate uint16 = 16000
)

// File is an audio file
type File struct {
	audioData WAV
}

//Writes the audio data to a file
func (f* File) WriteToFile (filename string) {
	data, err := audioData.Data()
	check(err)
	err = ioutil.WriteFile(filename, data, 0644)
	check(err)
	return
}

//Pads silence to both the left and right channels of the audio data
func (f* File) Pad (secs int16) {
	padding := make([]byte, secs * rate)
	oldData, err := f.audioData.Data()
	check(err)
	f.audioData, err := NewWAVFromData(padding)
	f.audioData.AddAudioData(oldData)
	f.audioData.AddAudioData(padding)
}

//Pads silence to the left channel of the audio data
func (f* File) PadLeft (secs int16) {
	padding := make([]byte, secs * rate)
	oldData, err := f.audioData.Data()
	check(err)
	f.audioData, err := NewWAVFromData(padding)
	f.audioData.AddAudioData(oldData)
}

//Pads silence to the right channel of the audio data
func (f* File) PadRight (secs int16) {
	padding := make([]byte, secs * rate)
	f.audioData.AddAudioData(padding)
}

//Trims all silence from the audio data
func (F* File) TrimSilence {
	audioData.TrimSilent(SilentThresh, 0)
}

//Plays the Audio File
func (f* File) Play {
	portaudio.initialize()
	stream, error := portaudio.OpenDefaultStream(NumChannels, NumChannels, Rate, )


	portaudio.terminate()
	return
}

//Creates a new Audio File from a recording
func NewFromRecording (length float64, silence_len float64) *File {
	return &File {
	}
}

//Creates a new Audio File from byte data
func NewFromBytes(b []byte) *File {
	// implement this
	var file File
	file.audioData, err := NewWAVFromData(b)
	check(err)
	return &file
}

//Creates a new Audio File from an io.Reader
func NewFromReader(r *io.Reader) *File {
	// implement this
	data, ioError := ioutil.ReadAll(r)
	check(ioError)
	wav, wavError := NewWAVFromData(data)
	check(wavError)
	return &File {
		audioData: wav
	}
	return nil
}

//Creates a new Audio File from an os.File
func NewFromFile(f *os.File) *File {
	return NewFromFileName(f.Name())
}

//Creates a new Audio File from a given filename
func NewFromFileName(f string) *File {
	data, err := ioutil.ReadFile(f)
	check(err)
	wav, err := NewWAVFromData(data)
	return &File {
		audioData: *wav
	}
}

//Creates a new Audio File from a Port Audio Stream
func NewFromStream(s portaudio.Stream) *File {
	return &File {

	}
}

//Returns the wav data contained in the audio file
func (f *File) WAVData() []byte {
	data, err = audioData.Data()
	check(err)
	return data
}

//Determines whether an audio slice is silent or not
func IsSilent(audio []byte) bool {
	var max byte = audio[0]
	for _, value := range array {
		if value > max {
			max = value
		}
	}
	return max < SilentThresh
} 
