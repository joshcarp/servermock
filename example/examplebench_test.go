//nolint: govet
package example

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"

	"github.com/joshcarp/servermock"
)

func BenchmarkServeRand(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	port, err := servermock.ServeRand(ctx, func(string, ...interface{}) {
	})
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		wg := &sync.WaitGroup{}
		for j := 0; j < 50; j++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				err = servermock.SetResponse(fmt.Sprintf("http://localhost:%d", port), servermock.Request{
					Path:       "/foo.service.bar.SomethingAPI/GetWhatever",
					Body:       []byte(`{"Hello": "true"}`),
					StatusCode: 200,
				})
				if err != nil {
					panic(err)
				}
				_, err := http.Get(fmt.Sprintf("http://localhost:%d/foo.service.bar.SomethingAPI/GetWhatever", port))
				if err != nil {
					panic(err)
				}
				wg.Done()
			}(wg)
		}
		wg.Wait()
	}
	cancel()
}
