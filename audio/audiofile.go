package audio

import (
	"bufio"
	"encoding/binary"
	"io"
	"io/ioutil"
	"os"

	"github.com/nkansal96/aurora-go/errors"

	"github.com/gordonklaus/portaudio"
)

const (
	BufSize      = 1 << 10
	SilentThresh = 1 << 10
	SampleRate   = 16000
	NumChannels  = 1
)

// File is an audio file
type File struct {
	AudioData *WAV

	playing    bool
	shouldStop bool
}

// Writes the audio data to a file
func (f *File) WriteToFile(filename string) error {
	return ioutil.WriteFile(filename, f.AudioData.Data(), 0644)
}

// Pad adds silence to both the beginning and end of the audio data. Silence
// is specified in seconds.
func (f *File) Pad(seconds float64) {
	// calculate number of bytes needed to pad given amount of seconds
	bytes := float64(f.AudioData.NumChannels*f.AudioData.BitsPerSample/8) * float64(f.AudioData.SampleRate) * seconds
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
	bytes := float64(f.AudioData.NumChannels*f.AudioData.BitsPerSample/8) * float64(f.AudioData.SampleRate) * seconds
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
	bytes := float64(f.AudioData.NumChannels*f.AudioData.BitsPerSample/8) * float64(f.AudioData.SampleRate) * seconds
	padding := make([]byte, int(bytes))

	// add the padding to the right
	f.AudioData.AddAudioData(padding)
}

// TrimSilence trims silence from both ends of the audio data
func (f *File) TrimSilence() {
	// TODO: calibrate constants
	f.AudioData.TrimSilent(0.03, 0.25)
}

func (f *File) Stop() {
	if f.playing {
		f.shouldStop = true
	}
}

// Play the audio file to the default output
func (f *File) Play() error {
	// initialize the underlying APIs for audio transmission
	portaudio.Initialize()
	defer portaudio.Terminate()

	// create a buffer for audio to be put into
	// hardcode 16-bit sample until i figure out a better way to do this
	bufLen := int(BufSize * f.AudioData.NumChannels)
	buf := make([]int16, bufLen)

	// create the audio stream to write to
	stream, err := portaudio.OpenDefaultStream(0, int(f.AudioData.NumChannels), float64(f.AudioData.SampleRate), BufSize, buf)
	if err != nil {
		return errors.NewFromErrorCodeInfo(errors.AudioFileNilStream, err.Error());
	}

	defer stream.Close()
	defer stream.Stop()
	stream.Start()

	// get audio data (without WAV header)
	data := f.AudioData.AudioData()
	step := bufLen * 2 // *2 since we need twice as many bytes to fill up `bufLen` in16s

	f.playing = true
	defer func() { f.playing = false }()

	for i := 0; i < len(data); i += step {
		// check if we should stop (user called f.Stop())
		if f.shouldStop {
			f.shouldStop = false
			break
		}
		// need to convert each 2-bytes in [i, i+step] to 1 little endian int16
		for j := 0; j < bufLen; j++ {
			k := j * 2
			buf[j] = int16(binary.LittleEndian.Uint16(data[i+k : i+k+2]))
		}
		// write the converted data into the stream
		err := stream.Write()
		if err != nil {
			return errors.NewFromErrorCodeInfo(errors.AudioFileNilStream, err.Error());
		}
	}

	return nil
}

// NewRecordingStream records audio data according to the given parameters (just
// like `NewFileFromRecording`), however instead of creating an *audio.File, it
// streams the data to a Reader, which can be read as soon as data is available.
// It emits a WAV header before the actual data, but the length of the data is
// not correct. Therefore, the resulting WAV should be read until EOF.
func NewRecordingStream(length float64, silenceLen float64) io.Reader {
	pr, pw := io.Pipe()
	bufwr := bufio.NewWriterSize(pw, 1000*BufSize)
	bufrd := bufio.NewReaderSize(pr, 1000*BufSize)

	go func() {
		defer pw.Close()
		defer bufwr.Flush()

		header := NewWAV().Data()
		bufwr.Write(header)

		ch := record(length, silenceLen)
		for d := range ch {
			if d.Error != nil {
				return
			}
			bufwr.Write(d.Data)
		}
	}()

	return bufrd
}

// NewFromRecording creates a new File by recording from the default input stream.
// length specifies the maximum length of the recording in seconds. silenceLen
// specifies how long in seconds to automatically stop after when consecutive
// silence is detected.
func NewFileFromRecording(length float64, silenceLen float64) (*File, error) {
	ch := record(length, silenceLen)
	audioData := make([]byte, 0)
	for d := range ch {
		if d.Error != nil {
			return nil, d.Error
		}

		audioData = append(audioData, d.Data...)
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
func NewFileFromBytes(b []byte) (*File, error) {
	wav, err := NewWAVFromData(b)
	if err != nil {
		return nil, err
	}
	return &File{wav, false, false}, err
}

//Creates a new Audio File from an io.Reader
func NewFileFromReader(r io.Reader) (*File, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return NewFileFromBytes(data)
}

//Creates a new Audio File from an os.File
func NewFileFromFile(f *os.File) (*File, error) {
	return NewFileFromReader(f)
}

//Creates a new Audio File from a given filename
func NewFileFromFileName(f string) (*File, error) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return NewFileFromBytes(data)
}

//Returns the wav data contained in the audio file
func (f *File) WAVData() []byte {
	return f.AudioData.Data()
}
