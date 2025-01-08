package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "http://localhost:8080"

type CreateTopicRequest struct {
	Name       string `json:"name"`
	BufferSize int    `json:"buffer_size"`
}

type PubMessage struct {
	Topic string `json:"topic"`
	Msg   string `json:"message"`
}

func main() {
	// Create a new topic
	createTopicReq := CreateTopicRequest{
		Name:       "test-topic",
		BufferSize: 10,
	}
	createTopic(createTopicReq)

	// Publish a message to the topic
	pubMessageReq := PubMessage{
		Topic: "test-topic",
		Msg:   "Hello, World!",
	}
	publishMessage(pubMessageReq)

	// Get the list of topics
	getTopics()

	// Get Messages from the topic "test-topic"
	getMessages("test-topic")
}

func getMessages(topic string) {

}

func createTopic(req CreateTopicRequest) {
	url := fmt.Sprintf("%s/topic", baseURL)
	jsonReq, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Error marshalling request:", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Create Topic Response:", string(body))
}

func publishMessage(req PubMessage) {
	url := fmt.Sprintf("%s/pub", baseURL)
	jsonReq, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Error marshalling request:", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Publish Message Response:", string(body))
}

func getTopics() {
	url := fmt.Sprintf("%s/topics", baseURL)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Get Topics Response:", string(body))
}
