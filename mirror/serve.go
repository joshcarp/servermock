package mirror

import (
	"golang.org/x/sync/errgroup"
	"net"
	"sync"
)

func Serve() {
	g := errgroup.Group{}
	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}
	sm := &sync.Map{}

	g.Go(GRPC(ln, sm))
	g.Go(HTTP(ln, sm))
	g.Wait()
}
