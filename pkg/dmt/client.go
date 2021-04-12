package dmt

import (
	"bytes"
	"encoding/json"
	"github.com/joshcarp/dmt/internal/data"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"net/http"
)

const setData = "Set-Data"

/* SetResponse sets the return to pth to that of b bytes */
func SetResponse(url string, pth string, body []byte, headers metadata.MD, isError bool, statusCode int) error {
	if len(pth) > 0 && pth[0] != '/' && pth != "/" {
		pth = "/" + pth
	}
	request := data.Request{
		Path:       pth,
		Headers:    headers,
		Body:       body,
		IsError:    isError,
		StatusCode: statusCode,
	}
	b, err := json.Marshal(request)
	if err != nil {
		return err
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
func SetGRPCResponse(url string, trace string, m proto.Message, headers metadata.MD, isError bool, statusCode int) error {
	b, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	return SetResponse(url, trace, b, headers, isError, statusCode)
}

/* SetGRPCResponse sets the return to pth to bytes marshaled from an interface */
func SetJsonResponse(url string, trace string, m interface{}, headers metadata.MD, isError bool, statusCode int) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return SetResponse(url, trace, b, headers, isError, statusCode)
}
