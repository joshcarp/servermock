package example

import (
	"context"
	"fmt"
	"github.com/joshcarp/dmt/pkg/dmt"
	"google.golang.org/grpc"
	"io"
	"net/http"
	"os"
)

func Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func ExampleSetResponse() {
	fmt.Println("------------------------")
	defer fmt.Println("------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	err := dmt.Serve(ctx, Printf, ":8000")
	if err != nil {
		panic(err)
	}
	err = dmt.SetResponse("http://localhost:8000", "/foo.service.bar.SomethingAPI/GetWhatever", []byte(`{"Hello": "true"}`), nil, false)
	if err != nil {
		panic(err)
	}
	resp, err := http.Get("http://localhost:8000/foo.service.bar.SomethingAPI/GetWhatever")
	if err != nil {
		panic(err)
	}
	_, _ = io.Copy(os.Stdout, resp.Body)
	fmt.Println()
	// or defer cancel()
	cancel()

	// Output:
	// ------------------------
	// Loading data for for request: /
	// Setting Data for request: /foo.service.bar.SomethingAPI/GetWhatever Length: 17
	// Loading data for for request: /foo.service.bar.SomethingAPI/GetWhatever
	// {"Hello": "true"}
	// ------------------------
}

func ExampleSetGRPCResponse() {
	fmt.Println("------------------------")
	defer fmt.Println("------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	err := dmt.Serve(ctx, Printf, ":8000")
	if err != nil {
		panic(err)
	}
	conn, err := grpc.Dial("localhost:8000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := NewExampleServiceClient(conn)

	err = dmt.SetGRPCResponse("http://localhost:8000", "/example.ExampleService/getExample", &Example{
		Name:     "ExampleName",
		Whatever: "ExampleFoo",
	}, nil, false)
	if err != nil {
		panic(err)
	}

	example, err := client.GetExample(context.Background(), &Example{})
	fmt.Println(example.Name)
	fmt.Println(example.Whatever)
	// or defer cancel()
	cancel()

	// Output:
	// ------------------------
	// Loading data for for request: /
	// Setting Data for request: /example.ExampleService/getExample Length: 25
	// Returning bytes for request: /example.ExampleService/getExample
	// ExampleName
	// ExampleFoo
	// ------------------------
}
