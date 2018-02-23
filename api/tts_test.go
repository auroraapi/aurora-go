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

func TestGetTTStEmptyString(t *testing.T) {
	_, err := api.GetTTS(c, "")
	require.NotNil(t, err)
	require.IsType(t, apiErrorType, err)
	require.Subset(t, []string{"APIInvalidInput"}, []string{err.(*errors.APIError).Code})
}
