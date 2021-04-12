package dmt

import (
	"encoding/json"
	"github.com/joshcarp/dmt/internal/data"
	"github.com/joshcarp/dmt/internal/unknown"
	"golang.org/x/sync/errgroup"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
)

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
		entry := data.Request{}
		if err = json.Unmarshal(b, &entry); err != nil {
			wr.WriteHeader(500)
			return
		}
		data.StoreData(s.sm, Endpoint, entry)
		s.log("Setting Data for request: %s Length: %d\n", Endpoint, len(entry.Body))
		return
	}
	s.log("Loading data for for request: %s\n", Endpoint)
	d, ok := data.LoadData(s.sm, Endpoint)
	if ok {
		wr.WriteHeader(d.StatusCode)
		for key, val := range d.Headers {
			wr.Header().Add(key, val[0])
		}
		_, _ = wr.Write(d.Body)
		return
	}
	wr.WriteHeader(500)
}

func (s server) ServeGRPC(_ interface{}, stream grpc.ServerStream) error {
	method, ok := grpc.MethodFromServerStream(stream)
	if !ok {
		return status.Error(codes.Unknown, "Unknown request")
	}
	s.log("Returning bytes for request: %s\n", method)
	d, ok := data.LoadData(s.sm, method)
	if !ok {
		return status.Error(codes.Unknown, "Unknown request")
	}
	err := stream.SetHeader(d.Headers)
	if err != nil {
		return err
	}
	if d.IsError {
		uproto := &spb.Status{}
		_ = proto.Unmarshal(d.Body, uproto)
		return status.FromProto(uproto).Err()
	}
	u := unknown.Unknown{}
	err = proto.Unmarshal(d.Body, &u)
	if err != nil {
		return err
	}
	return stream.SendMsg(&u)
}
