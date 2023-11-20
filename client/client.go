package client

import (
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

type (
	NodeAddress     = types.NodeAddress
	NodesList       = types.NodesList
	ConnectionsList = types.ConnectionsList
)

func StartClient(serverNode NodeAddress, bootstrapNode NodeAddress) {
	// Get the list of existing nodes in the network from the bootstrap node
	var existingNodes *NodesList = getExistingNodesInNetwork(bootstrapNode, serverNode)
	if existingNodes == nil {
		return
	}
	existingNodesList := *existingNodes

	// Connect to the existing nodes
	var nodeConnectionsList ConnectionsList = establishConnectionWithExistingNodes(existingNodesList)

	// // Start listening for incoming transactions
	// go listenForIncomingTransactions(serverNode)

	// // Start listening for incoming connections
	// go listenForIncomingConnections(serverNode)

	// // Start sending transactions to the existing nodes
	// go sendTransactionsToExistingNodes(nodeConnectionsList)

	// // Start sending connections to the existing nodes
	// go sendConnectionsToExistingNodes(nodeConnectionsList)

	// // Start sending transactions to the bootstrap node
	// go sendTransactionsToBootstrapNode(bootstrapNode)
}
