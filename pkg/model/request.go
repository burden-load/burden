package model

type Request struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Body    string            `json:"headers,omitempty"`
	Headers map[string]string `json:"body,omitempty"`
}
