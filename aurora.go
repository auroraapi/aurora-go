package aurora

import (
	"github.com/nkansal96/aurora-go/config"
)

// Config is an alias for config.C so that the user can easily
// override the SDK configuation by typing something like:
//
//   aurora.Config.AppID    = "My ID"
//   aurora.Config.AppToken = "My Token"
//
var Config = config.C
