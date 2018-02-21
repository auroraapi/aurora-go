package audio

import  (
	"os"
	"io"
	"math"
	"github.com/gordonklaus/portaudio"
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
	chk(err)
	err = ioutil.WriteFile(filename, data, 0644)
	chk(err)
	return
}

//Pads silence to both the left and right channels of the audio data
func (f* File) Pad (secs int16) {
	padding := make([]byte, secs * rate)
	oldData, err := f.audioData.Data()
	chk(err)
	f.audioData, err := NewWAVFromData(padding)
	chk(err)
	f.audioData.AddAudioData(oldData)
	f.audioData.AddAudioData(padding)
}

//Pads silence to the left channel of the audio data
func (f* File) PadLeft (secs int16) {
	padding := make([]byte, secs * rate)
	oldData, err := f.audioData.Data()
	chk(err)
	f.audioData, err := NewWAVFromData(padding)
	chk(err)
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
	defer portaudio.terminate()

	data, err := f.audioData.Data()
	chk(err)
	out := make([]byte, 64)
	stream, err := portaudio.OpenDefaultStream(0, NumChannels, Rate, len(out), out)
	chk(err)
	defer stream.Close()

	//Making the assumption that no audio file will be under 64 bytes long
	for i := 0; i <= len(data); i += 64 {
		copy(out, data[i:i+64])
		chk(stream.Write())
	}
	return
}

//Creates a new Audio File from a recording
//length is the length in seconds that the recording needs to be
//silence_len is the amount of time after which to stop recording if only silence is deteced
func NewFromRecording(length float64, silence_len float64) *File {
	portaudio.Initialize()
	defer portaudio.Terminate()

	data := ([]int16, 0)
	inBuffer := ([]int16, 64)

	stream, error := portaudio.OpenDefaultStream(NumChannels, 0, Rate, len(inBuffer), inBuffer)
	chk(error)
	defer stream.Close()

	chk(stream.Start())
	var silentFor float64 = 0
	for {
		chk(stream.Read())
		data = append(data, inBuffer)

		if IsSilent(inBuffer) {
			silentFor = silentFor + float64(len(inBuffer))/Rate
		}

		if length == 0 && silentFor > silence_len {
			break
		}

		if length > 0 && len(data) > int(length*Rate) {
			break
		}
	}

	var file File	
	file.audioData, err := NewWAVFromData(data)
	chk(err)
	return &file
}

//Creates a new Audio File from byte data
func NewFromBytes(b []byte) *File {
	// implement this
	var file File
	file.audioData, err := NewWAVFromData(b)
	chk(err)
	return &file
}

//Creates a new Audio File from an io.Reader
func NewFromReader(r *io.Reader) *File {
	// implement this
	data, ioError := ioutil.ReadAll(r)
	chk(ioError)
	wav, wavError := NewWAVFromData(data)
	chk(wavError)
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
	chk(err)
	wav, err := NewWAVFromData(data)
	chk(err)
	return &File {
		audioData: *wav
	}
}

//Creates a new Audio File from a Port Audio Stream
//Assumes that streams are opened and closed by the callee
//This function will terminate execution when the stream has no more bytes to send
//b is the buffer passed to the stream upon initialization
func NewFromStream(s portaudio.Stream, b []byte) *File {
	data := make([]byte, 0)
	numAvailableBytes, err := s.AvailableToRead()
	chk(err)
	while numAvailableBytes > 0 {
		chk(s.Read())
		data = append(data, b)
		numAvailableBytes, err := s.AvailableToRead()
	}

	var file File
	file.audioData, err = NewWAVFromData(data)
	return &file
}

//Returns the wav data contained in the audio file
func (f *File) WAVData() []byte {
	data, err = audioData.Data()
	chk(err)
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

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
