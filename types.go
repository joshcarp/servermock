package dmt

type Logger func(format string, args ...interface{})

type Request struct {
	Path       string              `json:"path"`
	Headers    map[string][]string `json:"headers"`
	Body       []byte              `json:"body"`
	IsError    bool                `json:"is_error"`
	StatusCode int                 `json:"status_code"`
	IsStack    bool                `json:"stack"`
	HeaderKeys HeaderPair          `json:"header_keys"`
}

type HeaderPair struct {
	Key string
	Val []string
}
