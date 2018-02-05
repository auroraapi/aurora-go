package backend

import (
	"io"
	"net/http"
	"net/url"
)

// MultipartFile is an in-memory representation of a file to upload
type MultipartFile struct {
	// Name is the form field name
	Name string
	// Filename is the file name to upload the file as
	FileName string
	// Data is the file data
	Data []byte
}

// CallParams describe the request to send
type CallParams struct {
	// Method is one of GET, POST, PATCH, DELETE, etc.
	Method string
	// Path is the relative path of the query
	Path string
	// Headers are any additional headers to send in the request
	Headers http.Header
	// Body is any data to send in the body
	Body io.Reader
	// Form is ignored for non-multipart calls. It is multipart form data
	Form url.Values
	// Query is querystring parameters
	Query url.Values
	// Files is ignored for non-multipart calls. It is multipart file data
	Files []MultipartFile
	// Credentials are the AppID, AppToken, etc. required to make the call
	Credentials *Credentials
}

// Credentials are the credentials for the API request
type Credentials struct {
	// AppID is the appliacation ID (sent for 'X-Application-ID' header)
	AppID string
	// AppToken is the application token (sent for 'X-Application-Token' header)
	AppToken string
	// DeviceID is the device ID (sent for 'X-Device-ID' header)
	DeviceID string
}

// Backend is an interface for a general backend that executes a given request
type Backend interface {
	// Set some properties of the backend
	SetBaseURL(url string)
	SetClient(client *http.Client)

	// Call the backend to perform a request
	Call(params *CallParams) (*http.Response, error)
	CallMultipart(params *CallParams) (*http.Response, error)

	// Lower level methods that can be called as well
	NewRequest(params *CallParams) (*http.Request, error)
	Do(req *http.Request) (*http.Response, error)
}
