package server

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net"
	"sync"
)

func GRPC(ln net.Listener, sm sync.Map) error {
	return grpc.NewServer(grpc.UnknownServiceHandler(func(srv interface{}, stream grpc.ServerStream) error {
		md, _ := metadata.FromIncomingContext(stream.Context())
		tracearr := md.Get("X-B3-TraceId")
		if len(tracearr) != 1 {
			return nil
		}
		e := empty.NewUnknown(data.LoadData(sm, tracearr[0]))
		stream.SendMsg(&e)
		return nil
	}),
	).Serve(ln)

}
