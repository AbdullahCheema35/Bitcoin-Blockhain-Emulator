package server

import (
	"log"
	"net"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/common"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

type (
	NodeAddress = types.NodeAddress
	Message     = types.Message
)

func receiveClientRequest(conn net.Conn) (bool, NodeAddress) {
	err, message := common.ReceiveMessage(conn)
	if !err {
		return false, types.NodeAddress{}
	}
	switch message.Header.Type {
	case types.MessageTypeConnectionRequest:
		sender := message.Header.Sender
		if message.Body != types.ConnectionRequestTypeSuccess {
			log.Printf("Client %v sent a connection request with non-success type\n", sender.GetAddress())
			return false, types.NodeAddress{}
		}
		log.Println("Received a connection request from", sender.GetAddress())
		return true, sender
	default:
		log.Println("Received an unknown message from", message.Header.Sender.GetAddress())
		return false, types.NodeAddress{}
	}
}

func sendResponseToClient(conn net.Conn, clientNodeAddress NodeAddress) bool {
	maxNeighbours := configuration.GetMaxNeighbours()
	currentNeighbours := configuration.LockCurrentNeighbours()
	currentConnections := configuration.LockCurrentConnections()
	defer func() {
		configuration.UnlockCurrentConnections(currentConnections)
		configuration.UnlockCurrentNeighbours(currentNeighbours)
	}()

	if currentNeighbours >= maxNeighbours {
		log.Println("Maximum neighbours reached")

		messageType := types.MessageTypeConnectionResponse
		sender := configuration.GetSelfServerAddress()
		messageHeader := types.NewMessageHeader(messageType, sender)
		messageBody := types.ConnectionResponseTypeFailure
		message := types.NewMessage(messageHeader, messageBody)

		common.SendMessage(conn, message)
		return false
	} else {
		messageType := types.MessageTypeConnectionResponse
		sender := configuration.GetSelfServerAddress()
		messageHeader := types.NewMessageHeader(messageType, sender)
		messageBody := types.ConnectionResponseTypeSuccess
		message := types.NewMessage(messageHeader, messageBody)

		common.SendMessage(conn, message)
		// increment the current neighbours
		currentNeighbours++
		// Add the client node address to the current connections
		clientNodeConnection := types.NewNodeConnection(clientNodeAddress, conn)
		currentConnections.AddNodeConnection(clientNodeConnection)
		log.Println("Current neighbours:", currentNeighbours)
		log.Println("Current connections:", currentConnections.GetNodeConnections())
		return true
	}
}

func respondToConnectionRequest(conn net.Conn) bool {
	var isRequestSuccess bool
	var clientNodeAddress NodeAddress

	isRequestSuccess, clientNodeAddress = receiveClientRequest(conn)
	if !isRequestSuccess {
		return false
	}

	isConnectionSuccess := sendResponseToClient(conn, clientNodeAddress)
	return isConnectionSuccess
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	log.Println("Received a connection request")

	var isConnectionSuccess bool = respondToConnectionRequest(conn)

	log.Println("Connection request handled; Success = ", isConnectionSuccess)

	if !isConnectionSuccess {
		return
	}

	// Now we can start listening for messages from the Client Node
	for {
		err, message := common.ReceiveMessage(conn)
		if !err {
			log.Println("Same issue")
			break
		}
		switch message.Header.Type {
		case types.MessageTypeTransaction:
			sender := message.Header.Sender
			// TODO: Handle the transaction
			// Temp fix
			transactionData := message.Body.(string)
			log.Printf("Received transaction %v from %v\n", transactionData, sender.GetAddress())
			// End of temp fix
		case types.MessageTypeBlock:
			sender := message.Header.Sender
			log.Println("Received a block from", sender.GetAddress())
		default:
			sender := message.Header.Sender
			log.Println("Received an unknown message from", sender.GetAddress())
		}
	}
}

func StartServer() {
	serverNode := configuration.GetSelfServerAddress()
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
