package audio_test

import (
	"os"
	"testing"

	"github.com/nkansal96/aurora-go/audio"
	"github.com/nkansal96/aurora-go/config"
	"github.com/nkansal96/aurora-go/errors"
	"github.com/stretchr/testify/require"
)

var apiErrorType *errors.APIError
var c *config.Config

func TestNewWAV(t *testing.T) {
	wav := audio.NewWAV()
	require.Equal(t, uint16(1), wav.NumChannels)
	require.Equal(t, uint32(16000), wav.SampleRate)
	require.Equal(t, uint16(1), wav.AudioFormat)
	require.Equal(t, uint16(16), wav.BitsPerSample)

	audioData := wav.AudioData()
	require.Equal(t, 0, len(audioData))
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
