package aurora

import (
	"github.com/nkansal96/aurora-go/api"
)

type Interpret struct {
	Intent   string
	Entities map[string]string
}

func NewInterpret(res *api.InterpretResponse) *Interpret {
	if res == nil {
		return nil
	}
	return &Interpret{Intent: res.Intent, Entities: res.Entities}
}
