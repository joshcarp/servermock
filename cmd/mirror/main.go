package main

import (
	"fmt"
	"github.com/joshcarp/mirror/mirror"
)

func main() {
	fmt.Println("Start server...")
	mirror.Serve(":8000")
}
