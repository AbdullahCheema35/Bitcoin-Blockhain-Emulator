package server

import (
	"encoding/gob"
	"log"
	"net"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

type (
	NodeAddress = types.NodeAddress
	Message     = types.Message
)

func respondToConnectionRequest(conn net.Conn) bool {
	maxNeighbours := configuration.MaxNeighbours
	currentNeighbours := <-configuration.CurrentNeighbours
	currentConnections := <-configuration.CurrentConnections
	defer func() {
		configuration.CurrentConnections <- currentConnections
		configuration.CurrentNeighbours <- currentNeighbours
	}()
	if currentNeighbours >= maxNeighbours {
		log.Println("Maximum neighbours reached")

		messageHeader := types.MessageTypeConnectionRequest
		messageBody := types.ConnectionRequestTypeFailure
		message := types.NewMessage(messageHeader, messageBody)

		enc := gob.NewEncoder(conn)
		err := enc.Encode(message)
		if err != nil {
			log.Println("Error encoding:", err)
		}
		return false
	} else {
		// increment the current neighbours
		currentNeighbours++

		log.Println("Current neighbours:", currentNeighbours)

		messageHeader := types.MessageTypeConnectionRequest
		messageBody := types.ConnectionRequestTypeFailure
		message := types.NewMessage(messageHeader, messageBody)

		enc := gob.NewEncoder(conn)
		err := enc.Encode(message)
		if err != nil {
			log.Println("Error encoding:", err)
		}
		// receive the response from the client
		dec := gob.NewDecoder(conn)
		var messageFromClient Message
		err = dec.Decode(&messageFromClient)
		if err != nil {
			log.Println("Error decoding:", err)
		}
		// check if the message is valid
		if messageFromClient.Header == types.MessageTypeConnectionResponse {
			// It means that client has sent its NodeAddress(Server's Address) in the body
			// So we can add it to our list of connections
			clientNodeAddress := messageFromClient.Body.(NodeAddress)

			// Create a new NodeConnection object
			clientNodeConnection := types.NewNodeConnection(clientNodeAddress, conn)

			// Add the connection to the list of connections
			currentConnections.AddNodeConnection(clientNodeConnection)

			log.Println("Connection established successfully with", clientNodeAddress.GetAddress())
		} else {
			log.Println("Client sent an invalid message")
		}
		return true
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	var isConnectionSuccess bool = respondToConnectionRequest(conn)

	if !isConnectionSuccess {
		return
	} else {
		// Now we can start listening for messages from the client
		dec := gob.NewDecoder(conn)
		for {
			var message Message
			err := dec.Decode(&message)
			if err != nil {
				log.Println("Error decoding:", err)
				break
			}
			switch message.Header {
			case types.MessageTypeTransaction:
				sender := message.Sender.(NodeAddress)
				log.Println("Received a transaction from", message.Sender.GetAddress())
			case types.MessageTypeBlock:
				sender := message.Sender.(NodeAddress)
				log.Println("Received a block from", sender.GetAddress())
			default:
				log.Println("Received an unknown message from", message.Sender.GetAddress())
			}
		}

	}
}

func StartServer(serverNode NodeAddress) {
	serverAddress := serverNode.GetAddress()
	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	log.Println("Server Node listening on port", serverNode.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}
