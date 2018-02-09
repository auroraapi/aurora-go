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
	name string
	audioData audio.Buffer

}

func (f* File) audioData *audio.Buffer {
	return &f.audioData
}

func NewFile (audio audio.Int16) *File {
	return &File {
		name: ""
		audioData: *NewBuffer(audio)
	}
}

func (f* File) WriteToFile (filename string) {

}

func (f* File) GetWav {

}

func (f* File) Pad (secs int16) {

}

func (f* File) PadLeft (secs int16) {

}

func (f* File) PadRight (secs int16) {

}

func (F* File) TrimSilence {
	return
}

func (f* File) Play {
	return
}


func NewFromRecording (length float64, silence_len float64) *File {
	return &File {
		name: "micdata"
	}
}

func NewFromBytes(b []byte) *File {
	// implement this
	var file *File {
		name: "data"
	}
	file.audioData() = *audio.NewBuffer(b)
	return file
}

func NewFromReader(r *io.Reader) *File {
	// implement this
	var audioData audio.Int16 = ReadFull 
	return nil
}

func NewFromFile(f *os.File) *File {
	return nil
}

func NewFromFileName(f string) *File {
	return &File {
		name: f
	}
}

func NewFromStream(s portaudio.Stream) *File {
	return &File {
		name: "streamdata"
	}
}

func (f *File) WAVData() []byte {
	// implement this
	return nil
}

func IsSilent(audio []byte) int16 {
	return -1
} 
