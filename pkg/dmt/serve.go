package dmt

import (
	"context"
	"github.com/googleapis/gax-go/v2"
	"golang.org/x/sync/errgroup"
	"net"
	"net/http"
	"sync"
	"time"
)

/* Serve servers a dmt server and blocks until the server is running. Use context.WithCancel to stop the server */
func Serve(ctx context.Context, addr string) error {
	g := errgroup.Group{}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	sm := &sync.Map{}
	g.Go(servegrpc(ln, sm))
	g.Go(servehttp(ln, sm))
	go func() {
		select {
		case <-ctx.Done():
			ln.Close()
		}
	}()
	bo := gax.Backoff{
		Initial:    time.Second,
		Multiplier: 2,
		Max:        10 * time.Second,
	}
	for {
		_, err := http.Get("http://localhost" + addr)
		if err != nil {
			if err := gax.Sleep(ctx, bo.Pause()); err != nil {

			}
			continue
		}
		return nil
	}
}
