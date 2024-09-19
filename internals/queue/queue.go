package queue

import (
	"fmt"
	"sync"
)

var BUFF_SIZE = 20

type Message struct {
	Content string
}

type Topic struct {
	Name    string
	Channel chan Message
}

type Exchange struct {
	Topics map[string]*Topic
	mu     sync.RWMutex
}

// Init Exchange
func NewExchange() *Exchange {
	return &Exchange{
		Topics: make(map[string]*Topic),
	}
}

// Create a Topic
func (e *Exchange) CreateTopic(name string, buffSize int) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.Topics[name]; !exists {
		e.Topics[name] = &Topic{
			Name:    name,
			Channel: make(chan Message, buffSize), // buffered Channels
		}
	}
}

// Publish Message to a Topic
func (e *Exchange) Publish(topicName string, msg Message) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if topic, exists := e.Topics[topicName]; exists {
		topic.Channel <- msg
	} else {
		fmt.Printf("Topic %s does not exist\n", topicName)
	}
}

// Return a Channel Object for requested Topic if exists
func (e *Exchange) Subscribe(topicName string) (<-chan Message, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if topic, exists := e.Topics[topicName]; exists {
		return topic.Channel, nil
	}
	return nil, fmt.Errorf("Topic %s does not exist", topicName)
}

// ListTopics returns a slice of topic names
func (e *Exchange) GetTopics() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var topics []string
	for name := range e.Topics {
		topics = append(topics, name)
	}
	return topics
}
