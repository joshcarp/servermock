package mirror

import (
	"fmt"
	"github.com/joshcarp/mirror/data"
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
	trace := r.Header.Get("X-B3-TraceId")
	if r.Header.Get("Set-Data") != "" {
		b, _ := ioutil.ReadAll(r.Body)
		data.StoreData(s.sm, trace, b)
		fmt.Printf("Setting Data for request: %s Lentgth: %d\n", trace, len(b))
		return
	}
	fmt.Printf("Loading data for for request: %s\n", trace)
	wr.Write(data.LoadData(s.sm, trace))
}
