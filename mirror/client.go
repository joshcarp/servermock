package mirror

import (
	"bytes"
	"encoding/json"
	"google.golang.org/protobuf/proto"
	"net/http"
)

const setData = "Set-Data"

func SetResponse(url string, pth string, b []byte) error {
	if len(pth) > 0 && pth[0] != '/' && pth != "/"{
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

func SetGRPCResponse(url string, trace string, m proto.Message) error {
	b, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	return SetResponse(url, trace, b)
}

func SetJsonResponse(url string, trace string, m interface{}) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return SetResponse(url, trace, b)
}
