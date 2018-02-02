package config

// Config configures the parameters that this SDK will use to operate
type Config struct {
	// AppID is the value to send as the `X-Application-ID` header
	AppID string

	// AppToken is the value to send as the `X-Application-Token` header
	AppToken string

	// DeviceID is the value to send as the `X-Device-ID` header (optional)
	DeviceID string
}

// C is an instance of the above config (with default values). It's exported
// so that all packages can use it
var C = &Config{}
