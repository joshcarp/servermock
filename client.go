package dmt

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

const HEADERMODE = "DMT-MODE"

/* SetResponse sets the return to pth to that of b bytes */
func SetResponse(url string, request Request) error {
	b, err := json.Marshal(request)
	if err != nil {
		return err
	}
	r, err := http.NewRequest(http.MethodPost, url+request.Path, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	r.Header.Add(HEADERMODE, "SET")
	_, err = http.DefaultClient.Do(r)
	return err
}

/* GetResponse Gets all of the requests that are currently stored for a key */
func GetResponses(url string, pth string) (reqs []Request, err error) {
	r, err := http.NewRequest(http.MethodGet, url+pth, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Add(HEADERMODE, "GET")
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &reqs); err != nil {
		return nil, err
	}
	return
}

/* SetGRPCResponse sets the return to pth to a bytes marshaled from a proto message */
func SetGRPCResponse(url string, m proto.Message, request Request) error {
	body, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	request.Body = body
	return SetResponse(url, request)
}

/* SetGRPCResponse sets the return to pth to bytes marshaled from an interface */
func SetJsonResponse(url string, pth string, m interface{}, headers metadata.MD, isError bool, statusCode int) error {
	body, err := json.Marshal(m)
	if err != nil {
		return err
	}
	request := Request{
		Path:       pth,
		Headers:    headers,
		Body:       body,
		IsError:    isError,
		StatusCode: statusCode,
	}
	return SetResponse(url, request)
}

func Reset(url string, pth string) error {
	r, err := http.NewRequest(http.MethodGet, url+pth, nil)
	if err != nil {
		return err
	}
	r.Header.Add(HEADERMODE, "RESET")
	_, err = http.DefaultClient.Do(r)
	return err
}
