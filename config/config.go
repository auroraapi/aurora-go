// Package config contains configuration information that is used by the
// rest of the SDK.
package config

import (
	"github.com/auroraapi/aurora-go/api/backend"
)

// Config configures the parameters that this SDK will use to operate
type Config struct {
	// AppID is the value to send as the `X-Application-ID` header
	AppID string

	// AppToken is the value to send as the `X-Application-Token` header
	AppToken string

	// DeviceID is the value to send as the `X-Device-ID` header (optional)
	DeviceID string

	// Backend to use for requests (configurable for testing purposes)
	Backend backend.Backend
}

// GetCredentials converts the client's credentials into a struct
// that gets passed into the backend.
func (c *Config) GetCredentials() *backend.Credentials {
	return &backend.Credentials{c.AppID, c.AppToken, c.DeviceID}
}

// C is an instance of the above config (with default values). It's exported
// so that all packages can use it.
var C = &Config{
	Backend: backend.NewAuroraBackend(),
}
