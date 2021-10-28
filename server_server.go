package servermock

import (
	"net/http"

	"google.golang.org/protobuf/types/known/emptypb"

	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func (s server) loadHttp(wr http.ResponseWriter, Endpoint string, f func(string) []string) error {
	s.log("Loading data for for request: %s\n", Endpoint)
	d, err := loadData(s.sm, Endpoint, f)
	if err != nil {
		wr.WriteHeader(500)
		return err

	}
	wr.WriteHeader(d.StatusCode)
	for key, val := range d.Headers {
		wr.Header().Add(key, val[0])
	}
	_, _ = wr.Write(d.Body)
	return nil
}

func (s server) ServeGRPC(_ interface{}, stream grpc.ServerStream) error {
	method, ok := grpc.MethodFromServerStream(stream)
	if !ok {
		return status.Error(codes.Unknown, "Unknown request")
	}
	md, _ := metadata.FromIncomingContext(stream.Context())
	s.log("Returning bytes for request: %s\n", method)
	entry, err := loadData(s.sm, method, func(s string) []string { return md.Get(s) })
	if err != nil {
		return err
	}
	if entry.IsError {
		uproto := &spb.Status{}
		_ = proto.Unmarshal(entry.Body, uproto)
		return status.FromProto(uproto).Err()
	}
	u := emptypb.Empty{}
	err = proto.Unmarshal(entry.Body, &u)
	if err != nil {
		return err
	}
	return stream.SendMsg(&u)
}
