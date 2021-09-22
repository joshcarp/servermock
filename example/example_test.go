//nolint: govet
package example

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joshcarp/dmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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

	err = dmt.SetResponse("http://localhost:8000", dmt.Request{
		Path:       "/foo.service.bar.SomethingAPI/GetWhatever",
		Body:       []byte(`{"Hello": "true"}`),
		StatusCode: 200,
	})
	if err != nil {
		panic(err)
	}
	resp, err := http.Get("http://localhost:8000/foo.service.bar.SomethingAPI/GetWhatever")
	if err != nil {
		panic(err)
	}
	_, _ = io.Copy(os.Stdout, resp.Body)
	fmt.Println()

	/* Make another request */
	resp, err = http.Get("http://localhost:8000/foo.service.bar.SomethingAPI/GetWhatever")
	if err != nil {
		panic(err)
	}
	_, _ = io.Copy(os.Stdout, resp.Body)
	fmt.Println()

	// or defer cancel()
	cancel()

	// Output:
	// ------------------------
	// Setting Data for request: /foo.service.bar.SomethingAPI/GetWhatever Length: 17
	// Loading data for for request: /foo.service.bar.SomethingAPI/GetWhatever
	// {"Hello": "true"}
	// Loading data for for request: /foo.service.bar.SomethingAPI/GetWhatever
	// {"Hello": "true"}
	// ------------------------
}

func ExampleSetResponseQueue() {
	fmt.Println("------------------------")
	defer fmt.Println("------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	err := dmt.Serve(ctx, Printf, ":8000")
	if err != nil {
		panic(err)
	}

	err = dmt.SetResponse("http://localhost:8000", dmt.Request{
		Path:       "/foo.service.bar.SomethingAPI/GetWhatever",
		Body:       []byte(`{"Hello": "true"}`),
		StatusCode: 200,
		IsQueue:    true,
	})
	if err != nil {
		panic(err)
	}
	err = dmt.SetResponse("http://localhost:8000", dmt.Request{
		Path:       "/foo.service.bar.SomethingAPI/GetWhatever",
		Body:       []byte(`{"Hello": "Blah"}`),
		StatusCode: 200,
		IsQueue:    true,
	})
	if err != nil {
		panic(err)
	}
	resp, err := http.Get("http://localhost:8000/foo.service.bar.SomethingAPI/GetWhatever")
	if err != nil {
		panic(err)
	}
	_, _ = io.Copy(os.Stdout, resp.Body)
	fmt.Println()

	/* Make another request */
	resp, err = http.Get("http://localhost:8000/foo.service.bar.SomethingAPI/GetWhatever")
	if err != nil {
		panic(err)
	}
	_, _ = io.Copy(os.Stdout, resp.Body)
	fmt.Println()

	// or defer cancel()
	cancel()

	// Output:
	// ------------------------
	// Setting Data for request: /foo.service.bar.SomethingAPI/GetWhatever Length: 17
	// Setting Data for request: /foo.service.bar.SomethingAPI/GetWhatever Length: 17
	// Loading data for for request: /foo.service.bar.SomethingAPI/GetWhatever
	// {"Hello": "true"}
	// Loading data for for request: /foo.service.bar.SomethingAPI/GetWhatever
	// {"Hello": "Blah"}
	// ------------------------
}

func ExampleReset() {
	fmt.Println("------------------------")
	defer fmt.Println("------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	err := dmt.Serve(ctx, Printf, ":8000")
	if err != nil {
		panic(err)
	}

	err = dmt.SetResponse("http://localhost:8000", dmt.Request{
		Path:       "/foo.service.bar.SomethingAPI/GetWhatever",
		Body:       []byte(`{"Hello": "true"}`),
		StatusCode: 200,
	})
	if err != nil {
		panic(err)
	}
	err = dmt.Reset("http://localhost:8000", "/foo.service.bar.SomethingAPI/GetWhatever")
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
	// Setting Data for request: /foo.service.bar.SomethingAPI/GetWhatever Length: 17
	// Loading data for for request: /foo.service.bar.SomethingAPI/GetWhatever
	// Error returning bytes: rpc error: code = Unknown desc = Unknown request
	// ------------------------
}

func ExampleSetGRPCResponse() {
	fmt.Println("------------------------")
	defer fmt.Println("------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	err := dmt.Serve(ctx, Printf, ":8001")
	if err != nil {
		panic(err)
	}
	conn, err := grpc.Dial("localhost:8001", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := NewExampleServiceClient(conn)

	err = dmt.SetGRPCResponse("http://localhost:8001", &Example{
		Name:     "ExampleName",
		Whatever: "ExampleFoo",
	}, dmt.Request{Path: "/example.ExampleService/getExample"})
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
		status.New(codes.Unknown, "Whatever123").Proto(),
		dmt.Request{
			Path:    "/example.ExampleService/getExample",
			Body:    nil,
			IsError: true,
		})
	if err != nil {
		panic(err)
	}

	_, err = client.GetExample(context.Background(), &Example{})
	fmt.Println(err)
	// or defer cancel()
	cancel()

	// Output:
	// ------------------------
	// Setting Data for request: /example.ExampleService/getExample Length: 15
	// Returning bytes for request: /example.ExampleService/getExample
	// rpc error: code = Unknown desc = Whatever123
	// ------------------------
}

func ExampleSetResponseError() {
	fmt.Println("------------------------")
	defer fmt.Println("------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	err := dmt.Serve(ctx, Printf, ":8000")
	if err != nil {
		panic(err)
	}
	err = dmt.SetResponse("http://localhost:8000", dmt.Request{
		Path:       "/foo.service.bar.SomethingAPI/GetWhatever",
		StatusCode: 404,
	})
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
	// Setting Data for request: /foo.service.bar.SomethingAPI/GetWhatever Length: 0
	// Loading data for for request: /foo.service.bar.SomethingAPI/GetWhatever
	// Error code: 404
	// ------------------------
}

func ExampleSetResponseHeaderKeys() {
	fmt.Println("------------------------")
	defer fmt.Println("------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	err := dmt.Serve(ctx, Printf, ":8001")
	if err != nil {
		panic(err)
	}
	conn, err := grpc.Dial("localhost:8001", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := NewExampleServiceClient(conn)

	err = dmt.SetGRPCResponse("http://localhost:8001", &Example{
		Name:     "ExampleName",
		Whatever: "ExampleFoo",
	}, dmt.Request{Path: "/example.ExampleService/getExample", HeaderKeys: dmt.HeaderPair{Key: "Authorisation", Val: []string{"Bearer foo123"}}})
	if err != nil {
		panic(err)
	}
	example, _ := client.GetExample(metadata.AppendToOutgoingContext(ctx, "Authorisation", "Bearer foo123"), &Example{})
	fmt.Println(example.Name)
	fmt.Println(example.Whatever)
	// or defer cancel()
	cancel()

	// Output:
	// ------------------------
	// Setting Data for request: /example.ExampleService/getExample Length: 25
	// Returning bytes for request: /example.ExampleService/getExample
	// ExampleName
	// ExampleFoo
	// ------------------------
}

func ExampleServeRand() {
	fmt.Println("------------------------")
	defer fmt.Println("------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	port, err := dmt.ServeRand(ctx, Printf)
	if err != nil {
		panic(err)
	}

	err = dmt.SetResponse(fmt.Sprintf("http://localhost:%d", port), dmt.Request{
		Path:       "/foo.service.bar.SomethingAPI/GetWhatever",
		Body:       []byte(`{"Hello": "true"}`),
		StatusCode: 200,
	})
	if err != nil {
		panic(err)
	}
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/foo.service.bar.SomethingAPI/GetWhatever", port))
	if err != nil {
		panic(err)
	}
	_, _ = io.Copy(os.Stdout, resp.Body)
	fmt.Println()

	/* Make another request */
	resp, err = http.Get(fmt.Sprintf("http://localhost:%d/foo.service.bar.SomethingAPI/GetWhatever", port))
	if err != nil {
		panic(err)
	}
	_, _ = io.Copy(os.Stdout, resp.Body)
	fmt.Println()

	// or defer cancel()
	cancel()

	// Output:
	// ------------------------
	// Setting Data for request: /foo.service.bar.SomethingAPI/GetWhatever Length: 17
	// Loading data for for request: /foo.service.bar.SomethingAPI/GetWhatever
	// {"Hello": "true"}
	// Loading data for for request: /foo.service.bar.SomethingAPI/GetWhatever
	// {"Hello": "true"}
	// ------------------------
}

func ExampleGetResponses() {
	fmt.Println("------------------------")
	defer fmt.Println("------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	port, err := dmt.ServeRand(ctx, Printf)
	if err != nil {
		panic(err)
	}

	err = dmt.SetResponse(fmt.Sprintf("http://localhost:%d", port), dmt.Request{
		Path:       "/foo.service.bar.SomethingAPI/GetWhatever",
		Body:       []byte(`{"Hello": "true"}`),
		StatusCode: 200,
	})
	if err != nil {
		panic(err)
	}

	reqs, err := dmt.GetResponses(fmt.Sprintf("http://localhost:%d", port), "/foo.service.bar.SomethingAPI/GetWhatever")
	if err != nil {
		panic(err)
	}
	fmt.Println(reqs)
	// or defer cancel()
	cancel()

	// Output:
	// ------------------------
	// Setting Data for request: /foo.service.bar.SomethingAPI/GetWhatever Length: 17
	// [{/foo.service.bar.SomethingAPI/GetWhatever map[] [123 34 72 101 108 108 111 34 58 32 34 116 114 117 101 34 125] false 200 false { []}}]
	// ------------------------
}

func ExampleSetGRPCQueue() {
	fmt.Println("------------------------")
	defer fmt.Println("------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := dmt.Serve(ctx, Printf, ":8001")
	if err != nil {
		panic(err)
	}
	conn, err := grpc.Dial("localhost:8001", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := NewExampleServiceClient(conn)

	err = dmt.SetGRPCResponse("http://localhost:8001", &Example{Name: "call 1"},
		dmt.Request{
			Path:    "/example.ExampleService/getExample",
			IsQueue: true,
		})
	if err != nil {
		panic(err)
	}
	err = dmt.SetGRPCResponse("http://localhost:8001", &Example{Name: "call 2"},
		dmt.Request{
			Path:    "/example.ExampleService/getExample",
			IsQueue: true,
		})
	if err != nil {
		panic(err)
	}
	example, _ := client.GetExample(context.Background(), &Example{})
	fmt.Println(example.Name)
	example, _ = client.GetExample(context.Background(), &Example{})
	fmt.Println(example.Name)
	// or defer cancel()
	cancel()

	// Output:
	// ------------------------
	// Setting Data for request: /example.ExampleService/getExample Length: 8
	// Setting Data for request: /example.ExampleService/getExample Length: 8
	// Returning bytes for request: /example.ExampleService/getExample
	// call 1
	// Returning bytes for request: /example.ExampleService/getExample
	// call 2
	// ------------------------
}
