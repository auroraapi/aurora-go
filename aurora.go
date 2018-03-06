// Package aurora is an SDK to interact with the Aurora API, making it
// easy to integrate voice user interfaces into your application.
package aurora

import (
	"github.com/nkansal96/aurora-go/config"
)

// Config is an alias for config.C so that you can easily
// override the SDK configuation by typing something like:
//
//   aurora.Config.AppID    = "My ID"
//   aurora.Config.AppToken = "My Token"
//   aurora.Config.DeviceID = "My Device"
var Config = config.C
