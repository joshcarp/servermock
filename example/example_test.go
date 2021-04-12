package example

import (
	"context"
	"fmt"
	"github.com/joshcarp/dmt/pkg/dmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	err = dmt.SetResponse("http://localhost:8000", "/foo.service.bar.SomethingAPI/GetWhatever", []byte(`{"Hello": "true"}`), nil, false, 200)
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
	}, nil, false, 200)
	if err != nil {
		panic(err)
	}

	example, _ := client.GetExample(context.Background(), &Example{})
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

func ExampleSetGRPCError() { //nolint: govet
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

	err = dmt.SetGRPCResponse("http://localhost:8000",
		"/example.ExampleService/getExample",
		status.New(codes.Unknown, "Whatever123").Proto(), nil, true, 0)
	if err != nil {
		panic(err)
	}

	_, err = client.GetExample(context.Background(), &Example{})
	fmt.Println(err)
	// or defer cancel()
	cancel()

	// Output:
	// ------------------------
	// Loading data for for request: /
	// Setting Data for request: /example.ExampleService/getExample Length: 15
	// Returning bytes for request: /example.ExampleService/getExample
	// rpc error: code = Unknown desc = Whatever123
	// ------------------------
}

func ExampleSetResponseError() { //nolint: govet
	fmt.Println("------------------------")
	defer fmt.Println("------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	err := dmt.Serve(ctx, Printf, ":8000")
	if err != nil {
		panic(err)
	}
	err = dmt.SetResponse("http://localhost:8000", "/foo.service.bar.SomethingAPI/GetWhatever",
		nil,
		nil,
		false,
		404)
	if err != nil {
		panic(err)
	}
	resp, err := http.Get("http://localhost:8000/foo.service.bar.SomethingAPI/GetWhatever")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Error code: %d\n", resp.StatusCode)
	// or defer cancel()
	cancel()

	// Output:
	// ------------------------
	// Loading data for for request: /
	// Setting Data for request: /foo.service.bar.SomethingAPI/GetWhatever Length: 0
	// Loading data for for request: /foo.service.bar.SomethingAPI/GetWhatever
	// Error code: 404
	// ------------------------
}
