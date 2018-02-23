package api_test

import (
	"testing"

	"github.com/nkansal96/aurora-go/api"
	"github.com/nkansal96/aurora-go/api/backend"
	"github.com/nkansal96/aurora-go/audio"
	"github.com/nkansal96/aurora-go/config"
	"github.com/nkansal96/aurora-go/errors"
	"github.com/nkansal96/aurora-go/testutils"
	"github.com/stretchr/testify/require"
)

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
	emptyWAVFile := testutils.CreateEmptyWAVFile()
	emptyAudioFile, _ := audio.NewFileFromBytes(emptyWAVFile)

	r, err := api.GetSTT(c, emptyAudioFile)
	require.Nil(t, err)
	require.NotNil(t, r)
	require.IsType(t, STTResponseType, r)
}
