package bootstrap

import (
	"log"
	"net"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/common"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

type NodeAddress = types.NodeAddress

type NodesList = types.NodesList

func addUpdateExistingNodesList(node NodeAddress) {
	_, currentExistingNodesMap := configuration.LockCurrentExistingNodes("")
	defer func() {
		configuration.UnlockCurrentExistingNodes(currentExistingNodesMap, "")
	}()
	currentExistingNodesMap.AddUpdateBootstrapNode(node)
}

func getExistingNodesList() NodesList {
	_, currentExistingNodesMap := configuration.LockCurrentExistingNodes("")
	defer func() {
		configuration.UnlockCurrentExistingNodes(currentExistingNodesMap, "")
	}()
	maxDelay := configuration.GetMaxSecondsPingDelay()
	existingNodesList := currentExistingNodesMap.GetRecentBootstrapNodes(maxDelay)
	return existingNodesList
}

func createMessageHeader(msgType types.MessageType) types.MessageHeader {
	sender := configuration.GetSelfBootstrapAddress()
	switch msgType {
	case types.MessageTypeBootstrapConnectionRequest:
		// log.Println("Received bootstrap connection request")
		msgType = types.MessageTypeBootstrapConnectionResponse
	case types.MessageTypeBootstrapPingRequest:
		// log.Println("Received bootstrap ping request")
		msgType = types.MessageTypeBootstrapPingResponse
	default:
		// log.Println("Invalid message type")
		msgType = types.MessageTypeUnknown
	}
	return types.MessageHeader{
		Type:   msgType,
		Sender: sender,
	}
}

func handleBootstrapRequest(msg types.Message) (types.MessageType, NodeAddress) {
	switch msg.Header.Type {
	case types.MessageTypeBootstrapConnectionRequest:
		sender := msg.Header.Sender
		// log.Println("Received bootstrap connection request from", sender.GetAddress())
		return types.MessageTypeBootstrapConnectionRequest, sender
	case types.MessageTypeBootstrapPingRequest:
		sender := msg.Header.Sender
		// log.Println("Received bootstrap ping request from", sender.GetAddress())
		return types.MessageTypeBootstrapPingRequest, sender
	default:
		// log.Println("Invalid message type")
		return types.MessageTypeUnknown, NodeAddress{}
	}
}

func respondToBootstrapRequest(conn net.Conn, msgType types.MessageType, currentExistingNodesList interface{}) {
	// log.Println("Sending bootstrap connection response")
	header := createMessageHeader(msgType)
	body := currentExistingNodesList
	msg := types.NewMessage(header, body)
	common.SendMessage(conn, msg)
}

// Handle the query from a node's server to the bootstrap node
// conn is the connection to the querying node's server
// nListChan is the read-only channel to receive the list of available nodes from handleAvailableNodesList goroutine
// nAddrChan is the write-only channel to send the querying node's server address to handleAvailableNodesList goroutine
func handleBootstrapQuery(conn net.Conn) {
	defer conn.Close()

	// Receive the NodeAddress of the node's server that is querying to the bootstrap node
	msgRcvSuccessfully, msg := common.ReceiveMessage(conn)
	if !msgRcvSuccessfully { // Failed to receive message from the querying node, i.e., the connection is broken
		// log.Println("Error receiving message from querying node")
		return
	}
	msgType, sender := handleBootstrapRequest(msg)
	switch msgType {
	case types.MessageTypeBootstrapConnectionRequest:
		// log.Println("Received bootstrap connection request")
		addUpdateExistingNodesList(sender)
		currentExistingNodesList := getExistingNodesList()
		respondToBootstrapRequest(conn, msgType, currentExistingNodesList)
	case types.MessageTypeBootstrapPingRequest:
		// log.Println("Received bootstrap ping request")
		addUpdateExistingNodesList(sender)
		respondToBootstrapRequest(conn, msgType, nil)
	default:
		// log.Println("Invalid message type")
		msgType = types.MessageTypeUnknown
		respondToBootstrapRequest(conn, msgType, nil)
	}
}

func StartBootstrapServer(bNode NodeAddress) {
	bootstrapAddress := bNode.GetAddress()
	listener, err := net.Listen("tcp", bootstrapAddress)
	if err != nil {
		log.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	log.Println("Bootstrap node listening on port", bNode.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go handleBootstrapQuery(conn)
	}
}
