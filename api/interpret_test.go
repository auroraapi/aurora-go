package api_test

import (
	"testing"
	"os"

	"github.com/nkansal96/aurora-go/api"
	"github.com/nkansal96/aurora-go/api/backend"
	"github.com/nkansal96/aurora-go/config"
	"github.com/nkansal96/aurora-go/errors"
	"github.com/stretchr/testify/require"
)

var apiErrorType *errors.APIError
var c *config.Config

func TestGetInterpretNoCredentials(t *testing.T) {
	badConfig := &config.Config{ Backend: backend.NewAuroraBackend() }

	_, err := api.GetInterpret(badConfig, "what is the weather in los angeles tomorrow")
	require.NotNil(t, err)
	require.IsType(t, apiErrorType, err)
	require.Subset(t, []string{"MissingApplicationID", "MissingApplicationToken"}, []string{err.(*errors.APIError).Code})
}

func TestGetInterpretEmptyString(t *testing.T) {
	_, err := api.GetInterpret(c, "")
	require.NotNil(t, err)
	require.IsType(t, apiErrorType, err)
	require.Subset(t, []string{"APIInvalidInput"}, []string{err.(*errors.APIError).Code})
}

func TestGetInterpret(t *testing.T) {
	r, err := api.GetInterpret(c, "remind me to pick up groceries tomorrow at 3 pm")
	require.Nil(t, err)
	require.NotNil(t, r)
	require.Equal(t, "set_reminder", r.Intent)
	require.Equal(t, "", r.Entities["location"])
	require.Equal(t, "pick up groceries", r.Entities["task"])
	require.Equal(t, "tomorrow at 3 pm", r.Entities["time"])
}

func TestMain(m *testing.M) {
	// set configuration from environment
	c = &config.Config{ 
		AppID: os.Getenv("APP_ID"),
		AppToken: os.Getenv("APP_TOKEN"),
		DeviceID: os.Getenv("DEVICE_ID"),
		Backend: backend.NewAuroraBackend(),
	}
	if len(os.Getenv("API_HOST")) > 0 {
		c.Backend.SetBaseURL(os.Getenv("API_HOST"))
	}

	// run tests
	os.Exit(m.Run())
}
