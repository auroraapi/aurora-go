package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/nkansal96/aurora-go/errors"
)

const (
	// the base URL for the backend (which is the actual service)
	baseURL = "https://api.auroraapi.com"
)

// AuroraBackend is an implementation that actually executes requests
// with the API server
type AuroraBackend struct {
	BaseURL string
	client  *http.Client
}

// NewAuroraBackend returns an AuroraBackend
func NewAuroraBackend() Backend {
	return &AuroraBackend{
		BaseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// NewAuroraBackendWithClient creates an AuroraBackend with the given client
func NewAuroraBackendWithClient(baseURL string, client *http.Client) Backend {
	return &AuroraBackend{BaseURL: baseURL, client: client}
}

// SetClient sets the http client for the backend
func (b *AuroraBackend) SetClient(client *http.Client) {
	b.client = client
}

// SetBaseURL sets the base URL for the backend
func (b *AuroraBackend) SetBaseURL(url string) {
	b.BaseURL = url
}

// Call implements a call to the backend
func (b *AuroraBackend) Call(params *CallParams) (*http.Response, error) {
	params.Path = fmt.Sprintf("%s%s?%s", b.BaseURL, params.Path, params.Query.Encode())
	req, err := b.NewRequest(params)
	if err != nil {
		return nil, err
	}

	return b.Do(req)
}

// CallMultipart implements a multipart call to the backend
func (b *AuroraBackend) CallMultipart(params *CallParams) (*http.Response, error) {
	// Pipe the output from the multipart writer so that we don't need to
	// buffer the file in memory before sending it over the network
	r, w := io.Pipe()

	// create a multipart form writing to the pipe
	multi := multipart.NewWriter(w)

	// don't block while copying data
	go func() {
		defer w.Close()
		defer multi.Close()
		// copy file data
		if params.Files != nil {
			for _, f := range params.Files {
				part, err := multi.CreateFormFile(f.Name, f.FileName)
				if err != nil {
					return
				}
				_, err = part.Write(f.Data)
				if err != nil {
					return
				}
			}
		}
		// copy form data
		if params.Form != nil {
			for k := range params.Form {
				multi.WriteField(k, params.Form.Get(k))
			}
		}
	}()

	// create the request
	params.Body = r
	params.Path = fmt.Sprintf("%s%s?%s", b.BaseURL, params.Path, params.Query.Encode())
	req, err := b.NewRequest(params)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-type", multi.FormDataContentType())
	return b.Do(req)
}

// NewRequest creates an http.Request from the given parameters
func (b *AuroraBackend) NewRequest(params *CallParams) (*http.Request, error) {
	req, err := http.NewRequest(params.Method, params.Path, params.Body)
	if err != nil {
		return nil, err
	}

	headers := params.Headers
	if headers == nil {
		headers = make(http.Header)
	}

	if params.Credentials != nil {
		headers.Add("X-Application-ID", params.Credentials.AppID)
		headers.Add("X-Application-Token", params.Credentials.AppToken)
		headers.Add("X-Device-ID", params.Credentials.DeviceID)
	}

	req.Header = headers
	return req, nil
}

// Do executes the given request
func (b *AuroraBackend) Do(req *http.Request) (*http.Response, error) {
	res, err := b.client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, handleError(res)
}

// handleError takes an http.Response object and assesses whether or not
// an error occurred. If it did, it returns an error object. Otherwise it
// returns nil
func handleError(r *http.Response) error {
	if r.StatusCode == http.StatusOK {
		return nil
	}

	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		var err errors.APIError
		json.NewDecoder(r.Body).Decode(&err)
		return &err
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return &errors.Error{
		Code:    "UnknownError",
		Message: string(body),
		Info:    fmt.Sprintf("An unexpected error occurred with code %d", r.StatusCode),
	}
}
