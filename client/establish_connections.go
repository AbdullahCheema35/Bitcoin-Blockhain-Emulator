package client

import (
	"log"
	"net"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/common"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

// Establishes a connection to the node at the given address and returns the connection or nil if the connection could not be established
func connectToNode(node NodeAddress) net.Conn {
	nodeAddress := node.GetAddress()
	conn, err := net.Dial("tcp", nodeAddress)
	if err != nil {
		log.Println("Error connecting to node:", err)
		return nil
	}
	return conn
}

func sendConnectionRequestToNode(node NodeAddress, conn net.Conn) bool {
	messageType := types.MessageTypeConnectionRequest
	sender := configuration.GetSelfServerAddress()
	messageHeader := types.NewMessageHeader(messageType, sender)
	messageBody := types.ConnectionRequestTypeSuccess
	message := types.NewMessage(messageHeader, messageBody)

	isMessageSent := common.SendMessage(conn, message)
	return isMessageSent
}

func receiveConnectionResponseFromNode(conn net.Conn) bool {
	isMessageReceived, message := common.ReceiveMessage(conn)
	if !isMessageReceived {
		return false
	}
	switch message.Header.Type {
	case types.MessageTypeConnectionResponse:
		if message.Body != types.ConnectionResponseTypeSuccess {
			log.Printf("Node %v sent a connection response with non-success type\n", message.Header.Sender.GetAddress())
			return false
		}
		log.Println("Received a successful connection response from", message.Header.Sender.GetAddress())
		return true
	default:
		log.Println("Received an unknown message from", message.Header.Sender.GetAddress())
		return false
	}
}

// Establishes connections with the nodes at the given addresses and returns the pointer to the connections' list
func establishConnectionWithExistingNodes(existingNodesList NodesList) {
	minNeighbours := configuration.GetMinNeighbours()
	currentNeighbours := configuration.LockCurrentNeighbours()
	currentConnections := configuration.LockCurrentConnections()
	defer func() {
		configuration.UnlockCurrentConnections(currentConnections)
		configuration.UnlockCurrentNeighbours(currentNeighbours)
	}()

	for _, node := range existingNodesList.GetNodes() {
		if currentNeighbours >= minNeighbours {
			break
		}
		var conn net.Conn = connectToNode(node)
		if conn != nil {
			isConnectionRequestSuccess := sendConnectionRequestToNode(node, conn)
			if isConnectionRequestSuccess {
				isConnectionResponseSuccess := receiveConnectionResponseFromNode(conn)
				if isConnectionResponseSuccess {
					currentNeighbours++
					nodeConnection := types.NewNodeConnection(node, conn)
					currentConnections.AddNodeConnection(nodeConnection)
					log.Println("Successfully Established connection with", node.GetAddress())
				}
			}
		}
	}
	log.Println("Established connections with", currentNeighbours, "neighbours")
}
