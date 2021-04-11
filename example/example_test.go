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

func ExampleSetResponse() {
	fmt.Println("------------------------")
	defer fmt.Println("------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := dmt.Serve(ctx, ":8000")
	if err != nil {
		log.Fatalf("Could not start server")
	}
	fmt.Println(dmt.SetResponse("http://localhost:8000", "/foo.service.bar.SomethingAPI/GetWhatever", []byte(`{"Hello": "true"}`)))
	resp, err := http.Get("http://localhost:8000/foo.service.bar.SomethingAPI/GetWhatever")
	if err != nil {
		panic(err)
	}
	_, _ = io.Copy(os.Stdout, resp.Body)
	fmt.Println()

	// Output:
	// ------------------------
	// Loading data for for request: /
	// Setting Data for request: /foo.service.bar.SomethingAPI/GetWhatever Length: 17
	// <nil>
	// Loading data for for request: /foo.service.bar.SomethingAPI/GetWhatever
	// {"Hello": "true"}
	// ------------------------
}

func ExampleSetGRPCResponse() {
	fmt.Println("------------------------")
	defer fmt.Println("------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := dmt.Serve(ctx, ":8000")
	if err != nil {
		log.Fatalf("Could not start server")
	}
	conn, err := grpc.Dial("localhost:8000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := NewExampleServiceClient(conn)

	err = dmt.SetGRPCResponse("http://localhost:8000", "/example.ExampleService/getExample", &Example{
		Name:     "ExampleName",
		Whatever: "ExampleFoo",
	})
	if err != nil {
		log.Fatalf("Could not set grpc response")
	}

	example, err := client.GetExample(context.Background(), &Example{})
	fmt.Println(err)
	fmt.Println(example.Name)
	fmt.Println(example.Whatever)

	// Output:
	// ------------------------
	// Loading data for for request: /
	// Setting Data for request: /example.ExampleService/getExample Length: 25
	// Returning bytes for request: /example.ExampleService/getExample
	// <nil>
	// ExampleName
	// ExampleFoo
	// ------------------------
}
