package servermock

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/soheilhy/cmux"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
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
	mux := cmux.New(ln)
	grpcL := mux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpL := mux.Match(cmux.Any())
	e.Go(mux.Serve)
	e.Go(func() error { return s.http.Serve(httpL) })
	e.Go(func() error { return s.grpc.Serve(grpcL) })
	return e.Wait()
}

func (s server) Stop() {
	s.grpc.Stop()
	s.http.Close()
}

func InterfaceToWriter(v interface{}, wr http.ResponseWriter) {
	b, err := json.Marshal(v)
	if err != nil {
		wr.WriteHeader(500)
	}
	_, _ = wr.Write(b)
}

func (s server) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	if r.ProtoMajor != 1 {
		return
	}
	Endpoint := r.URL.Path
	switch strings.ToLower(r.Header.Get(HEADERMODE)) {
	case "set":
		setData(wr, r, s.log, s.sm, Endpoint)
		return
	case "reset":
		resetData(s.sm)
		return
	case "get":
		reqs, ok := loadAllData(s.sm, Endpoint)
		if !ok {
			wr.WriteHeader(400)
		}
		InterfaceToWriter(reqs, wr)
		return
	default:
		err := s.loadHttp(wr, Endpoint, func(s string) []string { return r.Header.Values(s) })
		if err != nil {
			s.log("Error returning bytes: %v", err)
		}
	}
}
