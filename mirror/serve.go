package mirror

import (
	"golang.org/x/sync/errgroup"
	"net"
	"sync"
)

func Serve(addr string) {
	g := errgroup.Group{}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	sm := &sync.Map{}

	g.Go(GRPC(ln, sm))
	g.Go(HTTP(ln, sm))
	g.Wait()
}
