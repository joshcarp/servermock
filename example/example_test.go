package example

import (
	"context"
	"fmt"
	"github.com/joshcarp/dmt/pkg/dmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net/http"
	"os"
)

func ExampleServe() {
	fmt.Println("------------------------")
	defer fmt.Println("------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	dmt.Serve(ctx, ":8000")
	fmt.Println(dmt.SetResponse("http://localhost:8000", "/foo.service.bar.SomethingAPI/GetWhatever", []byte(`{"Hello": "true"}`)))
	resp, err := http.Get("http://localhost:8000/foo.service.bar.SomethingAPI/GetWhatever")
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, resp.Body)

	// Output:
	// Loading data for for request: /
	// Setting Data for request: /foo.service.bar.SomethingAPI/GetWhatever Length: 17
	// <nil>
	// Loading data for for request: /foo.service.bar.SomethingAPI/GetWhatever
	// {"Hello": "true"}
}

func ExampleGRPC() {
	fmt.Println("------------------------")
	defer fmt.Println("------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	dmt.Serve(ctx, ":8000")
	conn, err := grpc.Dial("localhost:8000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := NewExampleServiceClient(conn)

	dmt.SetGRPCResponse("http://localhost:8000", "/example.ExampleService/getExample", &Example{
		Name:     "ExampleName",
		Whatever: "ExampleFoo",
	})

	example, err := client.GetExample(context.Background(), &Example{})
	fmt.Println(err)
	fmt.Println(example)

	// Output:
	// ------------------------
	// Loading data for for request: /
	// Setting Data for request: /example.ExampleService/getExample Length: 25
	// Returning bytes for request: /example.ExampleService/getExample
	// <nil>
	// name:"ExampleName" whatever:"ExampleFoo"
	// ------------------------
}
