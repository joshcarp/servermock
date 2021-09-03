package dmt

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/googleapis/gax-go/v2"
)

/* Serve servers a dmt server and blocks until the server is running. Use context.WithCancel to stop the server */
func Serve(ctx context.Context, log Logger, addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return ServeLis(ctx, log, ln)
}

/* ServeLis servers a dmt server on a listener and blocks until the server is running. Use context.WithCancel to stop the server */
func ServeLis(ctx context.Context, log Logger, ln net.Listener) error {
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
		_, err := http.Get(fmt.Sprintf("http://localhost:%d", ln.Addr().(*net.TCPAddr).Port))
		if err != nil {
			if err := gax.Sleep(ctx, bo.Pause()); err != nil {
				return err
			}
			continue
		}
		return nil
	}
}

/* ServeRand servers a dmt server on a random port and blocks until the server is running. Use context.WithCancel to stop the server */
func ServeRand(ctx context.Context, log Logger) (int, error) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	return ln.Addr().(*net.TCPAddr).Port, ServeLis(ctx, log, ln)
}
