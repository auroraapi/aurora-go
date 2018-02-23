package api_test

import (
	"testing"

	"github.com/nkansal96/aurora-go/api"
	"github.com/nkansal96/aurora-go/api/backend"
	"github.com/nkansal96/aurora-go/config"
	"github.com/nkansal96/aurora-go/errors"
	"github.com/stretchr/testify/require"
)

func TestGetInterpretNoCredentials(t *testing.T) {
	badConfig := &config.Config{Backend: backend.NewAuroraBackend()}

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
	r, err := api.GetInterpret(c, "what time is it in los angeles")
	require.Nil(t, err)
	require.NotNil(t, r)
	require.Equal(t, "time", r.Intent)
	require.Equal(t, "los angeles", r.Entities["location"])
}
