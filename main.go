package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <port>")
		return
	}
	server := NewChatServer(":" + os.Args[1])
	server.Start()
}
