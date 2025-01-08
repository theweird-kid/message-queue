package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:7000")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server at localhost:8080")

	// Create a scanner to read user input
	scanner := bufio.NewScanner(os.Stdin)

	for {
		// Display the menu
		fmt.Println("\nMenu:")
		fmt.Println("1. SUBSCRIBE <topic>")
		fmt.Println("2. UNSUBSCRIBE <topic>")
		fmt.Println("3. PUBLISH <topic> <message>")
		fmt.Println("4. FETCH <topic>")
		fmt.Println("5. LIST_TOPICS")
		fmt.Println("6. CREATE_TOPIC <topic>")
		fmt.Println("7. EXIT")
		fmt.Print("Enter your choice: ")

		// Read user input
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		parts := strings.SplitN(input, " ", 3)

		if len(parts) < 1 {
			fmt.Println("Invalid input. Please try again.")
			continue
		}

		command := parts[0]

		switch command {
		case "SUBSCRIBE":
			if len(parts) < 2 {
				fmt.Println("Usage: SUBSCRIBE <topic>")
				continue
			}
			topic := parts[1]
			sendCommand(conn, fmt.Sprintf("SUBSCRIBE %s", topic))

		case "UNSUBSCRIBE":
			if len(parts) < 2 {
				fmt.Println("Usage: UNSUBSCRIBE <topic>")
				continue
			}
			topic := parts[1]
			sendCommand(conn, fmt.Sprintf("UNSUBSCRIBE %s", topic))

		case "PUBLISH":
			if len(parts) < 3 {
				fmt.Println("Usage: PUBLISH <topic> <message>")
				continue
			}
			topic := parts[1]
			message := parts[2]
			sendCommand(conn, fmt.Sprintf("PUBLISH %s %s", topic, message))

		case "FETCH":
			if len(parts) < 2 {
				fmt.Println("Usage: FETCH <topic>")
				continue
			}
			topic := parts[1]
			sendCommand(conn, fmt.Sprintf("FETCH %s", topic))

		case "LIST_TOPICS":
			sendCommand(conn, "LIST_TOPICS")

		case "CREATE_TOPIC":
			if len(parts) < 2 {
				fmt.Println("Usage: CREATE_TOPIC <topic>")
				continue
			}
			topic := parts[1]
			sendCommand(conn, fmt.Sprintf("CREATE_TOPIC %s", topic))

		case "EXIT":
			fmt.Println("Exiting...")
			return

		default:
			fmt.Println("Invalid command. Please try again.")
		}
	}
}

func sendCommand(conn net.Conn, command string) {
	// Send the command to the server
	fmt.Fprintf(conn, "%s\n", command)

	// Read the response from the server
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading response from server:", err)
		return
	}

	// Print the response
	fmt.Print("Server response: ", response)
}
