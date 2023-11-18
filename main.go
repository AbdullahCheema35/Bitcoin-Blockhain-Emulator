package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"strconv"
)

type Message struct {
	Content string
}

func handleConnection(conn net.Conn, messages chan Message) {
	defer conn.Close()

	decoder := gob.NewDecoder(conn)
	var msg Message

	err := decoder.Decode(&msg)
	if err != nil {
		fmt.Println("Error decoding message:", err)
		return
	}

	messages <- msg
}

func startServer(port int, messages chan Message) {
	address := ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Node listening on port", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn, messages)
	}
}

func sendMessageToNode(address string, msg Message) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	encoder := gob.NewEncoder(conn)
	err = encoder.Encode(msg)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run p2p.go <port>")
		return
	}

	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Invalid port number:", err)
		return
	}

	messages := make(chan Message)

	go startServer(port, messages)

	// Simulate sending a message to another node
	go func() {
		otherNodePort := 8081 // Change this to the port of another node
		otherNodeAddress := fmt.Sprintf("localhost:%d", otherNodePort)
		msg := Message{Content: "Hello from Node " + strconv.Itoa(port)}

		err := sendMessageToNode(otherNodeAddress, msg)
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
		fmt.Println("Message sent to Node", otherNodePort)
	}()

	// Receive messages from other nodes
	receivedMsg := <-messages
	fmt.Println("Received message:", receivedMsg.Content)
}
