package main

import (
	"fmt"

	"github.com/theweird-kid/message-queue/cmd/server"
)

func main() {
	fmt.Println("TCP Rewrite")
	server := server.NewServer("localhost", 7000)
	server.Start()
}
