package dmt

import (
	"github.com/joshcarp/dmt/internal/data"
	"github.com/joshcarp/dmt/internal/unknown"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
)

type Logger func(format string, args ...interface{})

type server struct {
	sm   *sync.Map
	http *http.Server
	grpc *grpc.Server
	log  Logger
}

func New(log Logger) *server {
	s := &server{sm: &sync.Map{}, log: log}
	s.http = &http.Server{Handler: s}
	s.grpc = grpc.NewServer(grpc.UnknownServiceHandler(s.ServeGRPC))
	return s

}

func (s server) Serve(ln net.Listener) error {
	e := errgroup.Group{}
	e.Go(func() error { return s.http.Serve(ln) })
	e.Go(func() error { return s.grpc.Serve(ln) })
	return e.Wait()
}

func (s server) Stop() {
	s.grpc.Stop()
	s.http.Close()
}

func (s server) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	if r.ProtoMajor != 1 {
		return
	}
	Endpoint := r.URL.Path
	if r.Header.Get("Set-Data") != "" {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.log("Error setting Data for request: %s\n", Endpoint)
		}
		data.StoreData(s.sm, Endpoint, b)
		s.log("Setting Data for request: %s Length: %d\n", Endpoint, len(b))
		return
	}
	s.log("Loading data for for request: %s\n", Endpoint)
	d := data.LoadData(s.sm, Endpoint)
	if len(d) != 0 {
		_, _ = wr.Write(d)
		return
	}
	wr.WriteHeader(500)
}

func (s server) ServeGRPC(srv interface{}, stream grpc.ServerStream) error {
	method, ok := grpc.MethodFromServerStream(stream)
	if !ok {
		return status.Error(codes.Unknown, "Unknown request")
	}
	s.log("Returning bytes for request: %s\n", method)
	u := unknown.Unknown{}
	err := proto.Unmarshal(data.LoadData(s.sm, method), &u)
	if err != nil {
		return err
	}
	return stream.SendMsg(&u)
}
