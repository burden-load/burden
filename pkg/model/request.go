package model

type Request struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Body    string            `json:"headers,omitempty"`
	Headers map[string]string `json:"body,omitempty"`
	Params  map[string]string `json:"params,omitempty"`
}

type PostmanCollection struct {
	Info PostmanInfo   `json:"info"`
	Item []PostmanItem `json:"item"`
}

type PostmanInfo struct {
	Name   string `json:"name"`
	Schema string `json:"schema"`
}

type PostmanItem struct {
	Name string           `json:"name"`
	Item []PostmanSubItem `json:"item"`
}

type PostmanSubItem struct {
	Name    string         `json:"name"`
	Request PostmanRequest `json:"request"`
}

type PostmanRequest struct {
	Method string        `json:"method"`
	Header []interface{} `json:"header"`
	URL    string        `json:"url"`
	Body   *PostmanBody  `json:"body,omitempty"`
}

type PostmanBody struct {
	Mode string `json:"mode"`
	Raw  string `json:"raw"`
}
