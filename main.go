package main

import (
	"fmt"
	"os"

	"github.com/Accelertor/Kurumirai/backend"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <port>")
		return
	}
	server := backend.NewServer(":" + os.Args[1])
	server.Start()
}
