package dmt

import (
	"context"
	"github.com/googleapis/gax-go/v2"
	"net"
	"net/http"
	"time"
)

/* Serve servers a dmt server and blocks until the server is running. Use context.WithCancel to stop the server */
func Serve(ctx context.Context, log Logger, addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s := New(log)
	go func() {
		_ = s.Serve(ln)
	}()
	go func() {
		<-ctx.Done()
		s.Stop()
		ln.Close()
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
				return err
			}
			continue
		}
		return nil
	}
}
