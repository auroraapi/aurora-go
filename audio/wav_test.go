package audio_test

import (
	"os"
	"bytes"
	"testing"
	"encoding/binary"

	"github.com/nkansal96/aurora-go/audio"
	"github.com/nkansal96/aurora-go/config"
	"github.com/nkansal96/aurora-go/errors"
	"github.com/stretchr/testify/require"
)

var apiErrorType *errors.APIError
var c *config.Config

func TestNewWAV(t *testing.T) {
	wav := audio.NewWAV()
	require.Equal(t, audio.DefaultNumChannels, wav.NumChannels)
	require.Equal(t, audio.DefaultSampleRate, wav.SampleRate)
	require.Equal(t, audio.DefaultAudioFormat, wav.AudioFormat)
	require.Equal(t, audio.DefaultBitsPerSample, wav.BitsPerSample)

	audioData := wav.AudioData()
	require.Equal(t, 0, len(audioData))
}

// If the WAVParams are specified, then the information should be 
// correctly written into the WAV struct
func TestNewWAVFromParamsCustom(t *testing.T) {
	emptyAudio := make([]byte, 0)
	wav := audio.NewWAVFromParams(&audio.WAVParams{1,10000,8,emptyAudio})
	require.Equal(t, uint16(1), wav.NumChannels)
	require.Equal(t, uint32(10000), wav.SampleRate)
	require.Equal(t, uint16(8), wav.BitsPerSample)
	require.Equal(t, 0, len(wav.AudioData()))
}

// If some of the WAVParams are specified to be 0, then the default
// parameters specified in the wav.go file will be given
func TestNewWAVFromParamsNotSpecified(t *testing.T) {
	emptyAudio := make([]byte, 0)
	wav := audio.NewWAVFromParams(&audio.WAVParams{1,0,16,emptyAudio})
	require.Equal(t, uint16(1), wav.NumChannels)
	require.Equal(t, uint32(16000), wav.SampleRate)
	require.Equal(t, uint16(16), wav.BitsPerSample)
	require.Equal(t, 0, len(wav.AudioData()))
}

func TestNewWAVFromData(t *testing.T) {
	emptyWAVFile := createEmptyWAVFile()

	wav, err := audio.NewWAVFromData(emptyWAVFile)
	require.Nil(t, err)
	require.Equal(t, uint16(1), wav.NumChannels)
	require.Equal(t, uint32(44100), wav.SampleRate)
	require.Equal(t, uint16(16), wav.BitsPerSample)
	require.Equal(t, 0, len(wav.AudioData()))
}


func TestNewWAVFromReader(t *testing.T) {
	emptyWAVFile := createEmptyWAVFile()
	r := bytes.NewReader(emptyWAVFile)

	wav, err := audio.NewWAVFromReader(r)
	require.Nil(t, err)
	require.Equal(t, uint16(1), wav.NumChannels)
	require.Equal(t, uint32(44100), wav.SampleRate)
	require.Equal(t, uint16(16), wav.BitsPerSample)
	require.Equal(t, 0, len(wav.AudioData()))
}

func TestAddAudioData(t *testing.T) {
	emptyWAVFile := createEmptyWAVFile()

	audioData := make([]byte, 4)
	binary.LittleEndian.PutUint32(audioData, 0x0000fdff)

	wav, err := audio.NewWAVFromData(emptyWAVFile)
	wav.AddAudioData(audioData)

	require.Nil(t, err)
	require.Equal(t, 4, len(wav.AudioData()))
}

func TestData(t *testing.T) {
	emptyWAVFile := createEmptyWAVFile()

	audioData := make([]byte, 4)
	binary.LittleEndian.PutUint32(audioData, 0x0000fdff)

	wav, err := audio.NewWAVFromData(emptyWAVFile)
	wav.AddAudioData(audioData)
	dataBytes := wav.Data() 

	require.Nil(t, err)
	require.Equal(t, []byte{0xff, 0xfd, 0x0, 0x0}, dataBytes[44:48])
}

func createEmptyWAVFile() []byte {
	emptyWAVFile := make([]byte, 44)
	// "RIFF" marker
	binary.BigEndian.PutUint32(emptyWAVFile[0:4], 0x52494646)
	// file size
	binary.LittleEndian.PutUint32(emptyWAVFile[4:8], 0)
	// "WAVE" type
	binary.BigEndian.PutUint32(emptyWAVFile[8:12], 0x57415645)
	// "fmt" section
	binary.BigEndian.PutUint32(emptyWAVFile[12:16], 0x666d7420)
	// length of fmt section
	binary.LittleEndian.PutUint32(emptyWAVFile[16:20], 16)
	// audio format
	binary.LittleEndian.PutUint16(emptyWAVFile[20:22], 1)
	// num channels
	binary.LittleEndian.PutUint16(emptyWAVFile[22:24], 1)
	// sample rate
	binary.LittleEndian.PutUint32(emptyWAVFile[24:28], 44100)
	// byte rate ((Sample Rate * Bit Size * Channels) / 8)
	binary.LittleEndian.PutUint32(emptyWAVFile[28:32], 44100)
	// block align ((bit size * channels) / 8)
	binary.LittleEndian.PutUint16(emptyWAVFile[32:34], 1)
	// bits per sample 
	binary.LittleEndian.PutUint16(emptyWAVFile[34:36], 16)
	// "data" marker
	binary.BigEndian.PutUint32(emptyWAVFile[36:40], 0x64617461)
	// data size 
	binary.LittleEndian.PutUint32(emptyWAVFile[40:44], 0)

	return emptyWAVFile
}

// TestMain sets up testing parameters and runs all tests
func TestMain(m *testing.M) {
	// set configuration from environment
	c = &config.Config{
		AppID:    os.Getenv("APP_ID"),
		AppToken: os.Getenv("APP_TOKEN"),
		DeviceID: os.Getenv("DEVICE_ID"),
	}
	if len(os.Getenv("API_HOST")) > 0 {
		c.Backend.SetBaseURL(os.Getenv("API_HOST"))
	}

	// run tests
	os.Exit(m.Run())
}
