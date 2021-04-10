package example

import (
	"context"
	"fmt"
	"github.com/joshcarp/dmt/dmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func ExampleServe() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go dmt.Serve(ctx, ":8000")
	time.Sleep(1)
	dmt.SetResponse("http://localhost:8000", "/foo.service.bar.SomethingAPI/GetWhatever", []byte("Hello"))
	resp, _ := http.Get("http://localhost:8000/foo.service.bar.SomethingAPI/GetWhatever")
	io.Copy(os.Stdout, resp.Body)
	// Output:
	// Setting Data for request: /foo.service.bar.SomethingAPI/GetWhatever Lentgth: 5
	// Loading data for for request: /foo.service.bar.SomethingAPI/GetWhatever
	// Hello
}

func ExampleGRPC(){
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go dmt.Serve(ctx, ":8000")
	time.Sleep(1)
	conn, err := grpc.Dial("localhost:8000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := NewExampleServiceClient(conn)

	dmt.SetGRPCResponse("http://localhost:8000", "/example.ExampleService/getExample", &Example{
		Name:          "ExampleName",
		Whatever:      "ExampleFoo",
	})

	example, err := client.GetExample(context.Background(), &Example{})
	fmt.Println(example)

	// Output:
	// Setting Data for request: /example.ExampleService/getExample Lentgth: 25
	// Returning bytes for request: /example.ExampleService/getExample
	// name:"ExampleName" whatever:"ExampleFoo"
}