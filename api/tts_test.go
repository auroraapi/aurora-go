package api_test

import (
	"testing"

	"github.com/nkansal96/aurora-go/api"
	"github.com/nkansal96/aurora-go/api/backend"
	"github.com/nkansal96/aurora-go/config"
	"github.com/nkansal96/aurora-go/errors"
	"github.com/stretchr/testify/require"
)

func TestGetTTSNoCredentials(t *testing.T) {
	badConfig := &config.Config{Backend: backend.NewAuroraBackend()}

	_, err := api.GetTTS(badConfig, "what is the weather in los angeles tomorrow")
	require.NotNil(t, err)
	require.IsType(t, apiErrorType, err)
	require.Subset(t, []string{"MissingApplicationID", "MissingApplicationToken"}, []string{err.(*errors.APIError).Code})
}

func TestGetTTSEmptyString(t *testing.T) {
	_, err := api.GetTTS(c, "")
	require.NotNil(t, err)
	require.IsType(t, apiErrorType, err)
	require.Subset(t, []string{"APIInvalidInput"}, []string{err.(*errors.APIError).Code})
}

func TestGetTTS(t *testing.T) {
	r, err := api.GetTTS(c, "what time is it in los angeles")
	require.Nil(t, err)
	require.NotNil(t, r)
	require.IsType(t, audioFileType, r)
}
