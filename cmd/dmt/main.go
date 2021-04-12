package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/joshcarp/dmt/pkg/dmt"
	"log"
)

func main() {
	p := flag.String("port", ":8000", "port to run dmt server on")
	flag.Parse()
	fmt.Println("Start server...")
	err := dmt.Serve(context.Background(), log.Printf, *p)
	if err != nil {
		log.Fatalf("Could not start server")
	}
	f := make(chan struct{})
	<-f
}
