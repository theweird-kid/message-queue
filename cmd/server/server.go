package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

type Message struct {
	Topic   string
	Content string
}

type Client struct {
	ID     string          // client id
	conn   net.Conn        // client connection
	topics map[string]bool // topics client is subscribed to
}

type Server struct {
	addr string // addr is the address of the server
	port int32  // port is the port number of the server

	clients  map[string]*Client    // clients is a map of client id to client
	topics   map[string][]*Client  // topics is a map of topic to clients
	messages map[string][]*Message // messages is a map of topic to messages

	mu sync.Mutex // mu is a mutex to protect the server state
}

func NewServer(addr string, port int32) *Server {
	srv := &Server{
		addr: addr,
		port: port,

		clients:  make(map[string]*Client),
		topics:   make(map[string][]*Client),
		messages: make(map[string][]*Message),
	}
	srv.loadState()
	return srv
}

func (s *Server) loadState() {

	// Load messages
	if data, err := os.ReadFile("./data/messages.json"); err == nil {
		json.Unmarshal(data, &s.messages)
	}

	// Load Topics
	if data, err := os.ReadFile("./data/topics.json"); err == nil {
		var topics map[string][]string
		json.Unmarshal(data, &topics)

		for topic, clientIDs := range topics {
			for _, clientID := range clientIDs {
				if client, exists := s.clients[clientID]; exists {
					client.topics[topic] = true
					s.topics[topic] = append(s.topics[topic], client)
				} else {
					client := &Client{
						ID: clientID, topics: make(map[string]bool),
					}

					s.clients[clientID] = client
					s.topics[topic] = append(s.topics[topic], client)
				}
			}
		}
	}
}

func (s *Server) saveState() {
	// Save messages
	if data, err := json.Marshal(s.messages); err == nil {
		os.WriteFile("./data/messages.json", data, 0644)
	}

	// Save Topics
	topics := make(map[string][]string)
	for topic, clients := range s.topics {
		for _, client := range clients {
			topics[topic] = append(topics[topic], client.ID)
		}
	}

	if data, err := json.Marshal(topics); err == nil {
		os.WriteFile("./data/topics.json", data, 0644)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Print when a connection is established
	fmt.Printf("New connection established from %s\n", conn.RemoteAddr().String())

	// Assign a unique ID to the client
	clientID := conn.RemoteAddr().String()
	client := &Client{
		ID: clientID, conn: conn, topics: make(map[string]bool),
	}

	// Save Client ID
	s.mu.Lock()
	s.clients[clientID] = client
	s.mu.Unlock()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("Received command: %s\n", line) // Log received command
		parts := strings.SplitN(line, " ", 3)
		if len(parts) < 1 {
			fmt.Fprintf(conn, "Invalid command\n")
			continue
		}

		command := parts[0]

		switch command {
		case "SUBSCRIBE":
			if len(parts) < 2 {
				fmt.Fprintf(conn, "Usage: SUBSCRIBE <topic>\n")
				continue
			}
			topic := parts[1]
			s.mu.Lock()
			client.topics[topic] = true
			s.topics[topic] = append(s.topics[topic], client)
			s.saveState()
			s.mu.Unlock()
			fmt.Fprintf(conn, "SUBSCRIBED to %s\n", topic)

		case "UNSUBSCRIBE":
			if len(parts) < 2 {
				fmt.Fprintf(conn, "Usage: UNSUBSCRIBE <topic>\n")
				continue
			}
			topic := parts[1]
			s.mu.Lock()
			if client.topics[topic] {
				delete(client.topics, topic)
				subscribers := s.topics[topic]
				for i, subscriber := range subscribers {
					if subscriber == client {
						s.topics[topic] = append(subscribers[:i], subscribers[i+1:]...)
						break
					}
				}
				s.saveState()
				fmt.Fprintf(conn, "UNSUBSCRIBED from %s\n", topic)
			} else {
				fmt.Fprintf(conn, "You are NOT subscribed to %s\n", topic)
			}
			s.mu.Unlock()

		case "PUBLISH":
			if len(parts) < 3 {
				fmt.Fprintf(conn, "Usage: PUBLISH <topic> <message>\n")
				continue
			}
			topic := parts[1]
			message := parts[2]
			s.mu.Lock()
			msg := Message{Topic: topic, Content: message}
			s.messages[topic] = append(s.messages[topic], &msg)
			s.saveState()
			for _, subscriber := range s.topics[topic] {
				if subscriber != client {
					fmt.Fprintf(subscriber.conn, "MESSAGE %s %s\n", topic, message)
				}
			}
			s.mu.Unlock()

		case "FETCH":
			if len(parts) < 2 {
				fmt.Fprintf(conn, "Usage: FETCH <topic>\n")
				continue
			}
			topic := parts[1]
			s.mu.Lock()
			if client.topics[topic] {
				for _, msg := range s.messages[topic] {
					fmt.Fprintf(conn, "MESSAGE %s %s\n", msg.Topic, msg.Content)
				}
			} else {
				fmt.Fprintf(conn, "You are NOT subscribed to %s\n", topic)
			}
			s.mu.Unlock()

		case "LIST_TOPICS":
			s.mu.Lock()
			fmt.Printf("Connection ClientID: %s has requested LIST_TOPICS\n", clientID)
			for topic := range s.topics {
				fmt.Fprintf(conn, "TOPIC %s\n", topic)
			}
			if len(s.topics) == 0 {
				fmt.Fprintf(conn, "No topics\n")
			}
			s.mu.Unlock()

		case "CREATE_TOPIC":
			if len(parts) < 2 {
				fmt.Fprintf(conn, "Usage: CREATE_TOPIC <topic>\n")
				continue
			}
			topic := parts[1]
			fmt.Printf("Connection ClientID: %s has requested CREATE_TOPIC %s\n", clientID, topic)
			s.mu.Lock()

			if _, exists := s.topics[topic]; exists {
				fmt.Fprintf(conn, "Topic %s already exists\n", topic)
			} else {
				s.topics[topic] = []*Client{}
				s.saveState()
				fmt.Fprintf(conn, "Topic %s created\n", topic)
			}
			fmt.Println("Done")
			fmt.Fprintf(conn, "Done\n")
			s.mu.Unlock()
		default:
			fmt.Fprintf(conn, "Unknown command: %s\n", command)
		}
	}

	s.mu.Lock()
	delete(s.clients, clientID)
	for topic := range client.topics {
		subscribers := s.topics[topic]
		for i, subscriber := range subscribers {
			if subscriber == client {
				s.topics[topic] = append(subscribers[:i], subscribers[i+1:]...)
				break
			}
		}
	}

	s.saveState()
	s.mu.Unlock()
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.addr, s.port))
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	fmt.Printf("Server started on %s:%d\n", s.addr, s.port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go s.handleConnection(conn)
	}
}
