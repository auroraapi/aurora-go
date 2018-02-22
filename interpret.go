package aurora

import (
	"github.com/nkansal96/aurora-go/api"
)

// Interpret contains the results from a call to the Aurora Interpret service.
type Interpret struct {
	// Intent represents the intent of the user. This can be an empty string if
	// the intent of the user was unclear. Otherwise, it will be one of the 
	// pre-determined values listed in the Aurora dashboard.
	Intent   string
	// Entities contain auxiliary information about the user's utterance. This
	// can be an empty map if no such information was detected. Otherwise, it
	// will be a key-value listing according to the entities described on the
	// Aurora dashboard.
	Entities map[string]string
}

// NewInterpret takes a response from the API and creates an Interpet object
// out of it. It doesn't really make sense for a developer to call this, but
// it is left exported in case it makes sense in the future.
func NewInterpret(res *api.InterpretResponse) *Interpret {
	if res == nil {
		return nil
	}
	return &Interpret{Intent: res.Intent, Entities: res.Entities}
}
