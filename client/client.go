package client

import (
	"log"
	"math/rand"
	"time"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/common"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

type (
	NodeAddress     = types.NodeAddress
	NodesList       = types.NodesList
	ConnectionsList = types.ConnectionsList
)

func sendArbitraryTransactionToNode(nodeConn types.NodeConnection) (bool, string) {
	conn := nodeConn.Conn
	messageType := types.MessageTypeTransaction
	sender := configuration.GetSelfServerAddress()
	messageHeader := types.NewMessageHeader(messageType, sender)
	// Transaction body is randomly constructed string for now
	// TODO: Make a transaction type

	// Create a random string of length 10
	rand.Seed(time.Now().UnixNano())

	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	length := 10

	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	randomString := string(b)
	// End of random string generation
	// Temp fix

	messageBody := randomString
	message := types.NewMessage(messageHeader, messageBody)

	isMessageSent := common.SendMessage(conn, message)

	return isMessageSent, randomString
}

func sendArbitraryTransactionToAllNodes(connectionsList ConnectionsList) {
	for _, nodeConn := range connectionsList.GetNodeConnections() {
		isMessageSent, transactionData := sendArbitraryTransactionToNode(nodeConn)
		if !isMessageSent {
			log.Printf("Could not send arbitrary transaction %v to %v\n", transactionData, nodeConn.Node.GetAddress())
		} else {
			log.Printf("Sent arbitrary transaction %v to %v\n", transactionData, nodeConn.Node.GetAddress())
		}
	}
}

func StartClient() {
	serverNode := configuration.GetSelfServerAddress()
	bootstrapNode := configuration.GetBootstrapNodeAddress()
	isSelfBootstrapNode := configuration.GetIsSelfBootstrapNode()

	if !isSelfBootstrapNode {
		// It is the bootstrap node of the network
		// Don't establish connection with the bootstrap node
		// Also, dont try to initialize the connections with the existing nodes in network, they will connect to this node, once they join the network

		// Get the list of existing nodes in the network from the bootstrap node
		var existingNodes *NodesList = getExistingNodesInNetwork(bootstrapNode, serverNode)
		if existingNodes == nil {
			log.Println("Could not get the list of existing nodes in the network. Exiting Client...")
			return
		}
		existingNodesList := *existingNodes

		log.Println("Received existing nodes in network. Length:", len(existingNodesList.GetNodes()))
		log.Println("Existing nodes in the network: ", existingNodesList.GetNodes())

		// Connect to the existing nodes
		establishConnectionWithExistingNodes(existingNodesList)
	}

	for {
		// Wait for ten seconds. Then send arbitrary transaction to all nodes
		time.Sleep(1000 * time.Millisecond)
		// Check if there are any neighbours
		currentNeighbours := configuration.ReadCurrentNeighbours()
		currentConnectionsList := configuration.ReadCurrentConnections()
		if currentNeighbours == 0 {
			log.Println("No neighbours to send arbitrary transaction to")
			continue
		}
		// Send arbitrary transaction to all nodes
		sendArbitraryTransactionToAllNodes(currentConnectionsList)
	}
}
