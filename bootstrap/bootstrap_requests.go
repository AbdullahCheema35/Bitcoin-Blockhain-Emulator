package bootstrap

import (
	"log"
	"net"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/comm"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

func getExistingNodesFromBootstrapNode(selfNode NodeAddress, conn net.Conn) interface{} {
	// create message header
	msgHeader := types.MessageHeader{
		Type:   types.MessageTypeBootstrapConnectionRequest,
		Sender: selfNode,
	}

	// Send the node's server address to the bootstrap node
	message := types.NewMessage(msgHeader, nil)
	success := comm.SendMessage(conn, message)

	if !success { // If the message was not sent successfully, return nil
		return nil
	}

	// Receive the list of available nodes from the bootstrap node
	success, msg := comm.ReceiveMessage(conn)
	if !success { // If the message was not received successfully, return nil
		return nil
	}

	// Handle the received message
	switch msg.Header.Type {
	case types.MessageTypeBootstrapConnectionResponse:
		// log.Println("Received bootstrap connection response")
		existingNodesList := msg.Body.(types.NodesList)
		existingNodesList.RemoveNode(selfNode)
		return existingNodesList
	default:
		log.Println("Invalid message type")
		return nil
	}
}

// Establishes a connection to the bootstrap server at the given address and returns the pointer to the connection
func connectToBootstrapNode(bootstrapNode NodeAddress) net.Conn {
	bootstrapAddress := bootstrapNode.GetAddress()
	conn, err := net.Dial("tcp", bootstrapAddress)
	if err != nil {
		log.Println("Couldn't connect to bootstrap node:", err)
		return nil
	} else {
		// log.Println("Connected to bootstrap node")
		return conn
	}
}

// returns the list of available nodes in the network from the bootstrap node
func GetExistingNodesInNetwork(bootstrapNode NodeAddress, selfNode NodeAddress) interface{} {
	// Connect to the bootstrap node
	bootstrapConn := connectToBootstrapNode(bootstrapNode)
	if bootstrapConn == nil {
		return nil
	}
	defer bootstrapConn.Close()

	// Get the list of available nodes from the bootstrap node
	existingNodesList := getExistingNodesFromBootstrapNode(selfNode, bootstrapConn)
	return existingNodesList
}

func PingBootstrapServer(bootstrapNode NodeAddress) bool {
	// Connect to the bootstrap node
	bootstrapConn := connectToBootstrapNode(bootstrapNode)
	if bootstrapConn == nil {
		return false
	}
	defer bootstrapConn.Close()

	// Get self node address
	sender := configuration.GetSelfServerAddress()

	// create message header
	msgHeader := types.MessageHeader{
		Type:   types.MessageTypeBootstrapPingRequest,
		Sender: sender,
	}

	// Send the node's server address to the bootstrap node
	message := types.NewMessage(msgHeader, nil)
	success := comm.SendMessage(bootstrapConn, message)

	if !success { // If the message was not sent successfully, return nil
		return false
	}

	// Receive the list of available nodes from the bootstrap node
	success, msg := comm.ReceiveMessage(bootstrapConn)
	if !success { // If the message was not received successfully, return nil
		return false
	}

	// Handle the received message
	switch msg.Header.Type {
	case types.MessageTypeBootstrapPingResponse:
		// log.Println("Received bootstrap connection response")
		return true
	default:
		log.Println("Invalid message type")
		return false
	}
}
