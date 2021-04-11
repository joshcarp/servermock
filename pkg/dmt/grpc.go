package dmt

import (
	"fmt"
	"github.com/joshcarp/dmt/internal/data"
	"github.com/joshcarp/dmt/internal/unknown"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"net"
	"sync"
)

func servegrpc(ln net.Listener, sm *sync.Map) func() error {
	return func() error {
		return grpc.NewServer(grpc.UnknownServiceHandler(func(srv interface{}, stream grpc.ServerStream) error {
			method, ok := grpc.MethodFromServerStream(stream)
			if !ok {
				return status.Error(codes.Unknown, "Unknown request")
			}
			fmt.Println("Returning bytes for request:", method)
			u := unknown.Unknown{}
			err := proto.Unmarshal(data.LoadData(sm, method), &u)
			if err != nil {
				return err
			}
			return stream.SendMsg(&u)
		}),
		).Serve(ln)
	}
}
