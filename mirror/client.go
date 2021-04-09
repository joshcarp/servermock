package mirror

import (
	"bytes"
	"google.golang.org/protobuf/proto"
	"net/http"
)

const traceID = "X-B3-TraceId"
const setData = "Set-Data"

func SetResponse(url string, trace string, b []byte) error {
	r, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	r.Header.Add(traceID, trace)
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
