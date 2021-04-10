package dmt

import (
	"fmt"
	"github.com/joshcarp/dmt/data"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
)

type Server struct {
	sm *sync.Map
}

func HTTP(ln net.Listener, sm *sync.Map) func() error {
	return func() error { return http.Serve(ln, Server{sm: sm}) }
}

func (s Server) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	if r.ProtoMajor != 1 {
		return
	}
	Endpoint := r.URL.Path
	if r.Header.Get("Set-Data") != "" {
		b, _ := ioutil.ReadAll(r.Body)
		data.StoreData(s.sm, Endpoint, b)
		fmt.Printf("Setting Data for request: %s Lentgth: %d\n", Endpoint, len(b))
		return
	}
	fmt.Printf("Loading data for for request: %s\n", Endpoint)
	wr.Write(data.LoadData(s.sm, Endpoint))
}
