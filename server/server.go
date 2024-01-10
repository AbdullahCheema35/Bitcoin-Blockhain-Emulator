package server

import (
	"log"
	"net"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/comm"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/connection"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/nodestate"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

type (
	NodeAddress     = types.NodeAddress
	Message         = types.Message
	ConnectionsList = types.ConnectionsList
)

func receiveClientRequest(nc types.NodeConnection) (bool, NodeAddress) {
	err, message := comm.ReceiveMessage(nc)
	if !err {
		return false, types.NodeAddress{}
	}
	switch message.Header.Type {
	case types.MessageTypeConnectionRequest:
		sender := message.Header.Sender
		// log.Println("Received a connection request from", sender.GetAddress())
		return true, sender
	default:
		// // log.Println("Received an unknown message from", message.Header.Sender.GetAddress())
		return false, types.NodeAddress{}
	}
}

func sendResponseToClient(nc types.NodeConnection) (bool, bool) {
	maxNeighbours := configuration.GetMaxNeighbours()
	currentNeighbours, currentConnections := nodestate.ReadCurrentConnections("")
	if currentNeighbours >= maxNeighbours || currentConnections.ExistsAddress(nc.Node) {
		log.Println("Maximum neighbours reached or client node already exists in the current connections list")

		messageType := types.MessageTypeFailure
		sender := configuration.GetSelfServerAddress()
		messageHeader := types.NewMessageHeader(messageType, sender)
		// messageBody := types.MessageTypeFailure
		message := types.NewMessage(messageHeader, nil)

		comm.SendMessage(nc, message)
		connectionSuccess, connectionClosed := false, false
		return connectionSuccess, connectionClosed
	} else {
		// Add the client node address to the current connections
		// clientNodeConnection := types.NewNodeConnection(clientNodeAddress, conn)
		_, currentConnections := nodestate.LockCurrentConnections("")
		success := connection.AddNewNodeConnection(&currentConnections, nc, "Server")
		nodestate.UnlockCurrentConnections(currentConnections, "")

		sender := configuration.GetSelfServerAddress()
		messageType := types.MessageTypeConnectionResponse
		if !success {
			messageType = types.MessageTypeFailure
		}
		messageHeader := types.NewMessageHeader(messageType, sender)
		message := types.NewMessage(messageHeader, nil)
		comm.SendMessage(nc, message)

		// log.Println("Current neighbours:", len(currentConnections.GetNodeConnections()))
		// log.Println("Current connections:", currentConnections.GetNodeConnections())
		connectionSuccess, connectionClosed := success, true
		return connectionSuccess, connectionClosed
	}
}

func respondToConnectionRequest(nc types.NodeConnection) (bool, bool) {
	var isRequestSuccess bool
	var clientNodeAddress NodeAddress

	isRequestSuccess, clientNodeAddress = receiveClientRequest(nc)
	if !isRequestSuccess {
		// log.Println("Unsuccessful connection request received from", clientNodeAddress.GetAddress())
		return false, false
	} else {
		// log.Println("Successful connection request received from", clientNodeAddress.GetAddress())
	}

	nc.SetNodeAddress(clientNodeAddress)

	isConnectionSuccess, isConnectionClosed := sendResponseToClient(nc)
	if !isConnectionSuccess {
		log.Println("Unsuccessful connection response sent to", clientNodeAddress.GetAddress())
	} else {
		// log.Println("Successful connection response sent to", clientNodeAddress.GetAddress())
	}
	return isConnectionSuccess, isConnectionClosed
}

func handleConnection(conn net.Conn) {
	// log.Println("---------------------------Received a connection request---------------------------------------")

	new_nc := types.NewNodeConnection(types.NewNodeAddress(0), conn)

	isConnectionSuccess, isConnectionClosed := respondToConnectionRequest(new_nc)

	if !isConnectionSuccess {
		if !isConnectionClosed {
			conn.Close()
		}
		return
	}

	if new_nc.Conn == nil {
		log.Panicln("Client Node Connection is nil, although everything went smooth")
	}

	// Now we can start listening for messages from the Client Node
	// connection.ListenForMessages(clientNodeConn.(types.NodeConnection))
}

func StartServer() {
	serverNode := configuration.GetSelfServerAddress()
	serverAddress := serverNode.GetAddress()
	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Panicln("Error listening:", err)
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
