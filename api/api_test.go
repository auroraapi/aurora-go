package api_test

import (
	"os"
	"testing"

	// "github.com/nkansal96/aurora-go/api"
	"github.com/nkansal96/aurora-go/api/backend"
	"github.com/nkansal96/aurora-go/config"
	"github.com/nkansal96/aurora-go/errors"
	// "github.com/stretchr/testify/require"
)

var apiErrorType *errors.APIError
var c *config.Config

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
