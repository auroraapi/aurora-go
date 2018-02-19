package backend_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/nkansal96/aurora-go/api/backend"
	"github.com/nkansal96/aurora-go/errors"
	"github.com/stretchr/testify/require"
)

func TestNewRequestEmptyCredentials(t *testing.T) {
	b := backend.NewAuroraBackend()
	r, err := b.NewRequest(&backend.CallParams{})
	require.Nil(t, err)
	require.Empty(t, r.Header, "Headers should be empty")
}

func TestNewRequestWithPartialCredentials(t *testing.T) {
	b := backend.NewAuroraBackend()
	r, err := b.NewRequest(&backend.CallParams{
		Credentials: &backend.Credentials{
			AppID:    "AppID",
			AppToken: "AppToken",
		},
	})

	require.Nil(t, err)
	require.Len(t, r.Header, 3)
	require.Empty(t, r.Header.Get("X-Device-ID"), "Device ID should be empty")
	require.Equal(t, "AppID", r.Header.Get("X-Application-ID"), "Mismatch App ID")
	require.Equal(t, "AppToken", r.Header.Get("X-Application-Token"), "Mismatch App Token")
}

func TestNewRequestWithFullCredentials(t *testing.T) {
	b := backend.NewAuroraBackend()
	r, err := b.NewRequest(&backend.CallParams{
		Credentials: &backend.Credentials{
			AppID:    "AppID",
			AppToken: "AppToken",
			DeviceID: "DeviceID",
		},
	})

	require.Nil(t, err)
	require.Len(t, r.Header, 3)
	require.Equal(t, "AppID", r.Header.Get("X-Application-ID"), "Mismatch App ID")
	require.Equal(t, "AppToken", r.Header.Get("X-Application-Token"), "Mismatch App Token")
	require.Equal(t, "DeviceID", r.Header.Get("X-Device-ID"), "Mismatch Device ID")
}

func TestNewRequestWithCustomHeaders(t *testing.T) {
	b := backend.NewAuroraBackend()
	h := make(http.Header)
	h.Add("X-Custom-Header-1", "CustomValue1")
	h.Add("X-Custom-Header-2", "CustomValue2")
	r, err := b.NewRequest(&backend.CallParams{Headers: h})

	require.Nil(t, err)
	require.Len(t, r.Header, 2)
	require.Equal(t, "CustomValue1", r.Header.Get("X-Custom-Header-1"), "Mismatch Custom Header 1")
	require.Equal(t, "CustomValue2", r.Header.Get("X-Custom-Header-2"), "Mismatch Custom Header 2")
}

func TestNewRequestWithCustomHeadersAndCredentials(t *testing.T) {
	b := backend.NewAuroraBackend()
	h := make(http.Header)
	h.Add("X-Custom-Header-1", "CustomValue1")
	h.Add("X-Custom-Header-2", "CustomValue2")
	r, err := b.NewRequest(&backend.CallParams{
		Headers: h,
		Credentials: &backend.Credentials{
			AppID:    "AppID",
			AppToken: "AppToken",
		},
	})

	require.Nil(t, err)
	require.Len(t, r.Header, 5)
	require.Empty(t, r.Header.Get("X-Device-ID"), "Device ID should be empty")
	require.Equal(t, "AppID", r.Header.Get("X-Application-ID"), "Mismatch App ID")
	require.Equal(t, "AppToken", r.Header.Get("X-Application-Token"), "Mismatch App Token")
	require.Equal(t, "CustomValue1", r.Header.Get("X-Custom-Header-1"), "Mismatch Custom Header 1")
	require.Equal(t, "CustomValue2", r.Header.Get("X-Custom-Header-2"), "Mismatch Custom Header 2")
}

func TestCallBadMethod(t *testing.T) {
	b := backend.NewAuroraBackend()
	_, err := b.Call(&backend.CallParams{Method: "BAD METHOD"})
	require.NotNil(t, err)
}

func TestCallBadPath(t *testing.T) {
	b := backend.NewAuroraBackend()
	_, err := b.Call(&backend.CallParams{Path: "af:::/_19bad-path"})
	require.NotNil(t, err)
}

func TestCallEmptyQueryString(t *testing.T) {
	eType := &errors.APIError{}
	b := backend.NewAuroraBackend()
	r, err := b.Call(&backend.CallParams{})
	require.Equal(t, "https", r.Request.URL.Scheme)
	require.Equal(t, "api.auroraapi.com", r.Request.URL.Host)
	require.Equal(t, "", r.Request.URL.Path)
	require.Equal(t, "", r.Request.URL.RawQuery)
	require.NotNil(t, err)
	require.IsType(t, eType, err)
	require.Equal(t, "NotFound", err.(*errors.APIError).Code)
}

func TestCallPathAndQuery(t *testing.T) {
	s := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(""))
	}))
	defer s.Close()

	b := backend.NewAuroraBackendWithClient(s.URL, s.Client())
	r, err := b.Call(&backend.CallParams{
		Path: "/v1/stt/",
		Query: url.Values(map[string][]string{
			"t1": []string{"v1"},
			"t2": []string{"v2"},
			"t3": []string{"v3", "v33"},
			"t4": []string{""},
		}),
	})
	require.Nil(t, err)
	require.Equal(t, "/v1/stt/", r.Request.URL.Path)
	require.Equal(t, "t1=v1&t2=v2&t3=v3&t3=v33&t4=", r.Request.URL.RawQuery)
}
