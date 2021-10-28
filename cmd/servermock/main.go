package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/joshcarp/servermock"
)

func main() {
	p := flag.String("port", ":8000", "port to run servermock server on")
	flag.Parse()
	fmt.Println("Start server...")
	err := servermock.Serve(context.Background(), log.Printf, *p)
	if err != nil {
		log.Fatalf("Could not start server")
	}
	f := make(chan struct{})
	<-f
}
