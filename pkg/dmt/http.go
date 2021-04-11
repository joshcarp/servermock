package dmt

import (
	"fmt"
	"github.com/joshcarp/dmt/internal/data"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
)

type server struct {
	sm *sync.Map
}

func servehttp(ln net.Listener, sm *sync.Map) func() error {
	return func() error { return http.Serve(ln, server{sm: sm}) }
}

func (s server) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	if r.ProtoMajor != 1 {
		return
	}
	Endpoint := r.URL.Path
	if r.Header.Get("Set-Data") != "" {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("Error setting Data for request: %s", Endpoint)
		}
		data.StoreData(s.sm, Endpoint, b)
		fmt.Printf("Setting Data for request: %s Length: %d\n", Endpoint, len(b))
		return
	}
	fmt.Printf("Loading data for for request: %s\n", Endpoint)
	d := data.LoadData(s.sm, Endpoint)
	if len(d) != 0 {
		_, _ = wr.Write(d)
		return
	}
	wr.WriteHeader(500)
}
