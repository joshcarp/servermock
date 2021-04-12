package data

type Request struct {
	Path       string              `json:"path"`
	Headers    map[string][]string `json:"headers"`
	Body       []byte              `json:"body"`
	IsError    bool                `json:"is_error"`
	StatusCode int               `json:"status_code"`
}
