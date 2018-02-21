package aurora_test

import (
	"os"
	"testing"

	aurora "github.com/nkansal96/aurora-go"
	"github.com/nkansal96/aurora-go/errors"
	"github.com/stretchr/testify/require"
)

var apiErrorType *errors.APIError

func TestTextInterpret(t *testing.T) {
	query := "what is the weather in los angeles"
	text := aurora.NewText(query)
	require.Equal(t, query, text.Text)

	i, err := text.Interpret()
	require.Nil(t, err)
	require.NotNil(t, i)
	require.Equal(t, "weather", i.Intent)
	require.Equal(t, "los angeles", i.Entities["location"])
}

func TestTextInterpretEmptyString(t *testing.T) {
	query := ""
	text := aurora.NewText(query)
	require.Equal(t, query, text.Text)

	_, err := text.Interpret()
	require.NotNil(t, err)
	require.IsType(t, apiErrorType, err)
}

// TestMain sets up testing parameters and runs all tests
func TestMain(m *testing.M) {
	// set configuration from environment
	aurora.Config.AppID = os.Getenv("APP_ID")
	aurora.Config.AppToken = os.Getenv("APP_TOKEN")
	aurora.Config.DeviceID = os.Getenv("DEVICE_ID")
	if len(os.Getenv("API_HOST")) > 0 {
		aurora.Config.Backend.SetBaseURL(os.Getenv("API_HOST"))
	}

	// run tests
	os.Exit(m.Run())
}
