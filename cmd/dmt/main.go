package main

import (
	"context"
	"fmt"
	"github.com/joshcarp/dmt/dmt"
)

func main() {
	fmt.Println("Start server...")
	dmt.Serve(context.Background(), ":8000")
}
