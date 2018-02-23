package api_test

import (
	"testing"
	"encoding/binary"

	"github.com/nkansal96/aurora-go/api"
	"github.com/nkansal96/aurora-go/api/backend"
	"github.com/nkansal96/aurora-go/audio"
	"github.com/nkansal96/aurora-go/config"
	"github.com/nkansal96/aurora-go/errors"
	"github.com/stretchr/testify/require"
)

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

func TestGetSTTNoCredentials(t *testing.T) {
	badConfig := &config.Config{Backend: backend.NewAuroraBackend()}

	_, err := api.GetTTS(badConfig, "what is the weather in los angeles tomorrow")
	require.NotNil(t, err)
	require.IsType(t, apiErrorType, err)
	require.Subset(t, []string{"MissingApplicationID", "MissingApplicationToken"}, []string{err.(*errors.APIError).Code})
}

// TODO: Not sure how to write a test for this
// func TestGetSTTIncorrectAudioFile(t *testing.T) {
// 	_, err := api.GetSTT(c, nil)
// 	require.NotNil(t, err)
// 	require.IsType(t, apiErrorType, err)
// 	require.Subset(t, []string{"APIInvalidInput"}, []string{err.(*errors.APIError).Code})
// }

func TestGetSTT(t *testing.T) {
	emptyWAVFile := createEmptyWAVFile()
	emptyAudioFile, _ := audio.NewFileFromBytes(emptyWAVFile)

	r, err := api.GetSTT(c, emptyAudioFile)
	require.Nil(t, err)
	require.NotNil(t, r)
	require.IsType(t, STTResponseType, r)
}
