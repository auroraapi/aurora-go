package api_test

import (
	"os"
	"testing"

	"github.com/auroraapi/aurora-go/api"
	"github.com/auroraapi/aurora-go/api/backend"
	"github.com/auroraapi/aurora-go/audio"
	"github.com/auroraapi/aurora-go/config"
	"github.com/auroraapi/aurora-go/errors"
)

var apiErrorType *errors.APIError
var c *config.Config
var audioFileType *audio.File
var STTResponseType *api.STTResponse

func TestMain(m *testing.M) {
	// set configuration from environment
	c = &config.Config{
		AppID:    os.Getenv("APP_ID"),
		AppToken: os.Getenv("APP_TOKEN"),
		DeviceID: os.Getenv("DEVICE_ID"),
		Backend:  backend.NewAuroraBackend(),
	}
	if len(os.Getenv("API_HOST")) > 0 {
		c.Backend.SetBaseURL(os.Getenv("API_HOST"))
	}

	// run tests
	os.Exit(m.Run())
}
