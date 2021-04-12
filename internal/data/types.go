package data

type Request struct {
	Path    string              `json:"path"`
	Headers map[string][]string `json:"headers"`
	Body    []byte              `json:"body"`
	IsError bool                `json:"is_error"`
}
