package dmt

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
	g.Go(servegrpc(ln, sm))
	g.Go(servehttp(ln, sm))
	select {
	case <-ctx.Done():
		return
	}

}
