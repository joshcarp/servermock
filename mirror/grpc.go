package mirror

import (
	"fmt"
	"github.com/joshcarp/mirror/data"
	"github.com/joshcarp/mirror/unknown"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"net"
	"sync"
)

func GRPC(ln net.Listener, sm *sync.Map) func() error {
	return func() error {
		return grpc.NewServer(grpc.UnknownServiceHandler(func(srv interface{}, stream grpc.ServerStream) error {
			md, _ := metadata.FromIncomingContext(stream.Context())
			tracearr := md.Get("X-B3-TraceId")
			if len(tracearr) != 1 {
				return nil
			}
			fmt.Println("Returning bytes for request: %s", tracearr[0])
			u := unknown.Unknown{}
			proto.Unmarshal(data.LoadData(sm, tracearr[0]), &u)
			return stream.SendMsg(&u)
		}),
		).Serve(ln)
	}
}
