package dmt

import (
	"bytes"
	"encoding/json"
	"google.golang.org/protobuf/proto"
	"net/http"
)

const setData = "Set-Data"

/* SetResponse sets the return to pth to that of b bytes */
func SetResponse(url string, pth string, b []byte) error {
	if len(pth) > 0 && pth[0] != '/' && pth != "/" {
		pth = "/" + pth
	}
	r, err := http.NewRequest(http.MethodPost, url+pth, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	r.Header.Add(setData, "t")
	_, err = http.DefaultClient.Do(r)
	return err
}

/* SetGRPCResponse sets the return to pth to a bytes marshaled from a proto message */
func SetGRPCResponse(url string, trace string, m proto.Message) error {
	b, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	return SetResponse(url, trace, b)
}

/* SetGRPCResponse sets the return to pth to bytes marshaled from an interface */
func SetJsonResponse(url string, trace string, m interface{}) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return SetResponse(url, trace, b)
}
