package dmt

import (
	"net/http"
	"reflect"

	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func (s server) loadHttp(wr http.ResponseWriter, Endpoint string) {
	s.log("Loading data for for request: %s\n", Endpoint)
	d, ok := loadData(s.sm, Endpoint)
	if ok {
		wr.WriteHeader(d.StatusCode)
		for key, val := range d.Headers {
			wr.Header().Add(key, val[0])
		}
		_, _ = wr.Write(d.Body)
		return
	}
	wr.WriteHeader(500)
}

func (s server) ServeGRPC(_ interface{}, stream grpc.ServerStream) error {
	method, ok := grpc.MethodFromServerStream(stream)
	if !ok {
		return status.Error(codes.Unknown, "Unknown request")
	}
	md, _ := metadata.FromIncomingContext(stream.Context())
	s.log("Returning bytes for request: %s\n", method)
	d, ok := loadAllData(s.sm, method)
	if !ok {
		return status.Error(codes.Unknown, "Unknown request")
	}
	var entry Request
	for _, dd := range d {
		mdd := md.Get(dd.HeaderKeys.Key)
		if reflect.DeepEqual(mdd, dd.HeaderKeys.Val) {
			entry = dd
			break
		}
	}
	if entry.IsError {
		uproto := &spb.Status{}
		_ = proto.Unmarshal(entry.Body, uproto)
		return status.FromProto(uproto).Err()
	}
	u := anypb.Any{}
	err := proto.Unmarshal(entry.Body, &u)
	if err != nil {
		return err
	}
	return stream.SendMsg(&u)
}
