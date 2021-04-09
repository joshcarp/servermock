package main

import (
	"context"
	"fmt"
	"github.com/joshcarp/mirror/mirror"
)

func main() {
	fmt.Println("Start server...")
	mirror.Serve(context.Background(), ":8000")
}
