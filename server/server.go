package server

import (
	"io/ioutil"
	"net"
	"net/http"
	"sync"
)

type Server struct {
	sm sync.Map
}

func ServeHTTP(ln net.Listener, sm sync.Map) error {
	return http.Serve(ln, Server{sm: sm})
}

func (s Server) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	if r.ProtoMajor != 1 {
		return
	}
	trace := r.Header.Get("X-B3-TraceId")
	if r.Header.Get("Set-Data") != "" {
		b, _ := ioutil.ReadAll(r.Body)
		data.StoreData(s.sm, trace, b)
		return
	}
}
