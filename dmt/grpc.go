package dmt

import (
	"fmt"
	"github.com/joshcarp/dmt/data"
	"github.com/joshcarp/dmt/unknown"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"net"
	"sync"
)

func GRPC(ln net.Listener, sm *sync.Map) func() error {
	return func() error {
		return grpc.NewServer(grpc.UnknownServiceHandler(func(srv interface{}, stream grpc.ServerStream) error {
			method, ok := grpc.MethodFromServerStream(stream)
			if !ok {
				return nil
			}
			fmt.Println("Returning bytes for request:", method)
			u := unknown.Unknown{}
			proto.Unmarshal(data.LoadData(sm, method), &u)
			return stream.SendMsg(&u)
		}),
		).Serve(ln)
	}
}
