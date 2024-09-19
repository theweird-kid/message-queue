package main

import (
	"fmt"

	"github.com/theweird-kid/message-queue/cmd/api"
	"github.com/theweird-kid/message-queue/internals/queue"
)

func main() {
	fmt.Println("Message Queue")

	// Instantitate Exchange and Topics
	e := queue.NewExchange()
	// topics
	e.CreateTopic("topic_1", queue.BUFF_SIZE)
	e.CreateTopic("topic_2", queue.BUFF_SIZE)

	topics := e.GetTopics()
	for _, topic := range topics {
		fmt.Println(topic)
	}

	// Instantiate Server Object
	s := api.NewServer(":8080", e)
	// Register Routes
	s.HandleRoutes()
	// Run the Server
	s.Run()
}
