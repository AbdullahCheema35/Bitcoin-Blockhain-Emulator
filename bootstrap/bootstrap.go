package bootstrap

import (
	"net"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/comm"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/nodestate"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

type NodeAddress = types.NodeAddress

type NodesList = types.NodesList

func addUpdateExistingNodesList(node NodeAddress) {
	_, currentExistingNodesMap := nodestate.LockCurrentExistingNodes("")
	defer func() {
		nodestate.UnlockCurrentExistingNodes(currentExistingNodesMap, "")
	}()
	currentExistingNodesMap.AddUpdateBootstrapNode(node)
}

func getExistingNodesList() NodesList {
	_, currentExistingNodesMap := nodestate.LockCurrentExistingNodes("")
	defer func() {
		nodestate.UnlockCurrentExistingNodes(currentExistingNodesMap, "")
	}()
	maxDelay := configuration.GetMaxSecondsPingDelay()
	existingNodesList := currentExistingNodesMap.GetRecentBootstrapNodes(maxDelay)
	return existingNodesList
}

func createMessageHeader(msgType types.MessageType) types.MessageHeader {
	sender := configuration.GetSelfBootstrapAddress()
	switch msgType {
	case types.MessageTypeBootstrapConnectionRequest:
		// // log.Println("Received bootstrap connection request")
		msgType = types.MessageTypeBootstrapConnectionResponse
	case types.MessageTypeBootstrapPingRequest:
		// // log.Println("Received bootstrap ping request")
		msgType = types.MessageTypeBootstrapPingResponse
	default:
		// // log.Println("Invalid message type")
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
		// // log.Println("Received bootstrap connection request from", sender.GetAddress())
		return types.MessageTypeBootstrapConnectionRequest, sender
	case types.MessageTypeBootstrapPingRequest:
		sender := msg.Header.Sender
		// // log.Println("Received bootstrap ping request from", sender.GetAddress())
		return types.MessageTypeBootstrapPingRequest, sender
	default:
		// // log.Println("Invalid message type")
		return types.MessageTypeUnknown, NodeAddress{}
	}
}

func respondToBootstrapRequest(nodeConn types.NodeConnection, msgType types.MessageType, currentExistingNodesList interface{}) {
	// // log.Println("Sending bootstrap connection response")
	header := createMessageHeader(msgType)
	body := currentExistingNodesList
	msg := types.NewMessage(header, body)
	comm.SendMessage(nodeConn, msg)
}

// Handle the query from a node's server to the bootstrap node
// conn is the connection to the querying node's server
// nListChan is the read-only channel to receive the list of available nodes from handleAvailableNodesList goroutine
// nAddrChan is the write-only channel to send the querying node's server address to handleAvailableNodesList goroutine
func handleBootstrapQuery(conn net.Conn) {
	defer conn.Close()

	// Create a new nodeConnection that will be used for communication
	nodeConn := types.NewNodeConnection(types.NewNodeAddress(0), conn)

	// Receive the NodeAddress of the node's server that is querying to the bootstrap node
	msgRcvSuccessfully, msg := comm.ReceiveMessage(nodeConn)
	if !msgRcvSuccessfully { // Failed to receive message from the querying node, i.e., the connection is broken
		// // log.Println("Error receiving message from querying node")
		return
	}
	msgType, sender := handleBootstrapRequest(msg)

	// Update nodeConn to include the NodeAddress of the querying node's server
	nodeConn.SetNodeAddress(sender)

	switch msgType {
	case types.MessageTypeBootstrapConnectionRequest:
		// // log.Println("Received bootstrap connection request")
		addUpdateExistingNodesList(sender)
		currentExistingNodesList := getExistingNodesList()
		respondToBootstrapRequest(nodeConn, msgType, currentExistingNodesList)
	case types.MessageTypeBootstrapPingRequest:
		// // log.Println("Received bootstrap ping request")
		addUpdateExistingNodesList(sender)
		respondToBootstrapRequest(nodeConn, msgType, nil)
	default:
		// // log.Println("Invalid message type")
		msgType = types.MessageTypeUnknown
		respondToBootstrapRequest(nodeConn, msgType, nil)
	}
}

func StartBootstrapServer(bNode NodeAddress) {
	bootstrapAddress := bNode.GetAddress()
	listener, err := net.Listen("tcp", bootstrapAddress)
	if err != nil {
		// log.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	// log.Println("Bootstrap node listening on port", bNode.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			// log.Println("Error accepting connection:", err)
			continue
		}
		go handleBootstrapQuery(conn)
	}
}
