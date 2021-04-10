package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/joshcarp/dmt/pkg/dmt"
)

func main() {
	p := flag.String("port", ":8000", "port to run dmt server on")
	flag.Parse()
	fmt.Println("Start server...")
	dmt.Serve(context.Background(), *p)
}
