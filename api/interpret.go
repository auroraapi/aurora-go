package api

type InterpretResponse struct {
	// Text is the original query
	Text string `json:"text"`

	Intent string `json:"intent"`

	Entities map[string]string `json:"entities"`
}

func GetInterpret(text string) (*InterpretResponse, error) {
	return nil, nil
}
