package mirror

import (
	"context"
	"golang.org/x/sync/errgroup"
	"net"
	"sync"
)

func Serve(ctx context.Context, addr string) {
	g := errgroup.Group{}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	sm := &sync.Map{}
	g.Go(GRPC(ln, sm))
	g.Go(HTTP(ln, sm))
	select {
	case <-ctx.Done():
		return
	}

}
