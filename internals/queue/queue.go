package queue

import "fmt"

var BUFF_SIZE int = 20

type Message struct {
	Content string
}

type Topic struct {
	Name    string
	Channel chan Message
}

type Exchange struct {
	Topics map[string]*Topic
}

// Init Exchange
func NewExchange() *Exchange {
	return &Exchange{
		Topics: make(map[string]*Topic),
	}
}

// Create a Topic
func (e *Exchange) CreateTopic(name string, buffSize int) {
	if _, exists := e.Topics[name]; !exists {
		e.Topics[name] = &Topic{
			Name:    name,
			Channel: make(chan Message, buffSize), // buffered Channels
		}
	}
}

// Publish Message to a Topic
func (e *Exchange) Publish(topicName string, msg Message) {
	if topic, exists := e.Topics[topicName]; exists {
		topic.Channel <- msg
	} else {
		fmt.Printf("Topic %s does not exist\n", topicName)
	}
}

// Return a Channel Object for requested Topic if exists
func (e *Exchange) Subscribe(topicName string) (<-chan Message, error) {
	if topic, exists := e.Topics[topicName]; exists {
		return topic.Channel, nil
	}
	return nil, fmt.Errorf("Topic %s does not exist", topicName)
}
