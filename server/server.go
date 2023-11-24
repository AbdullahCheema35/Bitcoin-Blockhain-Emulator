package server

import (
	"log"
	"net"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/comm"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/connection"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

type (
	NodeAddress     = types.NodeAddress
	Message         = types.Message
	ConnectionsList = types.ConnectionsList
)

func receiveClientRequest(conn net.Conn) (bool, NodeAddress) {
	err, message := comm.ReceiveMessage(conn)
	if !err {
		return false, types.NodeAddress{}
	}
	switch message.Header.Type {
	case types.MessageTypeConnectionRequest:
		sender := message.Header.Sender
		// log.Println("Received a connection request from", sender.GetAddress())
		return true, sender
	default:
		// log.Println("Received an unknown message from", message.Header.Sender.GetAddress())
		return false, types.NodeAddress{}
	}
}

func sendResponseToClient(conn net.Conn, clientNodeAddress NodeAddress) (bool, bool) {
	maxNeighbours := configuration.GetMaxNeighbours()
	currentNeighbours, currentConnections := configuration.ReadCurrentConnections("")
	if currentNeighbours >= maxNeighbours || currentConnections.ExistsAddress(clientNodeAddress) {
		// log.Println("Maximum neighbours reached or client node already exists in the current connections list")

		messageType := types.MessageTypeFailure
		sender := configuration.GetSelfServerAddress()
		messageHeader := types.NewMessageHeader(messageType, sender)
		// messageBody := types.MessageTypeFailure
		message := types.NewMessage(messageHeader, nil)

		comm.SendMessage(conn, message)
		connectionSuccess, connectionClosed := false, false
		return connectionSuccess, connectionClosed
	} else {
		// Add the client node address to the current connections
		clientNodeConnection := types.NewNodeConnection(clientNodeAddress, conn)
		_, currentConnections := configuration.LockCurrentConnections("")
		success := connection.AddNewNodeConnection(&currentConnections, clientNodeConnection, "Server")
		configuration.UnlockCurrentConnections(currentConnections, "")

		sender := configuration.GetSelfServerAddress()
		messageType := types.MessageTypeConnectionResponse
		if !success {
			messageType = types.MessageTypeFailure
		}
		messageHeader := types.NewMessageHeader(messageType, sender)
		message := types.NewMessage(messageHeader, nil)
		comm.SendMessage(conn, message)

		// log.Println("Current neighbours:", len(currentConnections.GetNodeConnections()))
		// log.Println("Current connections:", currentConnections.GetNodeConnections())
		connectionSuccess, connectionClosed := success, true
		return connectionSuccess, connectionClosed
	}
}

func respondToConnectionRequest(conn net.Conn) (bool, bool) {
	var isRequestSuccess bool
	var clientNodeAddress NodeAddress

	isRequestSuccess, clientNodeAddress = receiveClientRequest(conn)
	if !isRequestSuccess {
		log.Println("Unsuccessful connection request received from", clientNodeAddress.GetAddress())
		return false, false
	} else {
		log.Println("Successful connection request received from", clientNodeAddress.GetAddress())
	}

	isConnectionSuccess, isConnectionClosed := sendResponseToClient(conn, clientNodeAddress)
	if !isConnectionSuccess {
		log.Println("Unsuccessful connection response sent to", clientNodeAddress.GetAddress())
	} else {
		log.Println("Successful connection response sent to", clientNodeAddress.GetAddress())
	}
	return isConnectionSuccess, isConnectionClosed
}

func handleConnection(conn net.Conn) {
	// log.Println("Received a connection request")

	_, isConnectionClosed := respondToConnectionRequest(conn)

	if !isConnectionClosed {
		log.Println("Closing connection. Line 97 of server.go")
		conn.Close()
	}
	// Now we can start listening for messages from the Client Node
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
