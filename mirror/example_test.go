package mirror

import (
	"io"
	"net/http"
	"os"
	"time"
)

func ExampleServe() {
	go Serve(":8000")
	time.Sleep(1)
	SetResponse("http://localhost:8000", "/foo.service.bar.SomethingAPI/GetWhatever", []byte("Hello"))
	resp, _ := http.Get("http://localhost:8000/foo.service.bar.SomethingAPI/GetWhatever")
	io.Copy(os.Stdout, resp.Body)

	// Output:
	// Setting Data for request: /foo.service.bar.SomethingAPI/GetWhatever Lentgth: 5
	// Loading data for for request: /foo.service.bar.SomethingAPI/GetWhatever
	// Hello
}
