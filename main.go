package main

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"net"
	"sync"
)

func main() {
	fmt.Println("Start server...")
	g := errgroup.Group{}
	// listen on port 8000
	ln, _ := net.Listen("tcp", ":8000")
	var sm sync.Map

	g.Go(server.GRPC(ln, sm))
	g.Go(server.ServeHTTP(ln, sm))
	g.Wait()

}
