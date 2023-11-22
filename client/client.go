package client

import (
	"log"
	"time"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/connection"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/propagation"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

type (
	NodeAddress     = types.NodeAddress
	NodesList       = types.NodesList
	ConnectionsList = types.ConnectionsList
)

func StartClient() {
	// Connect to the network in safe mode
	connection.ConnectWithNetwork_SafeMode()

	// Normal execution of the client
	for {
		// Wait for ten seconds. Then send arbitrary transaction to all nodes
		// log.Println("Waiting for one seconds")
		time.Sleep(30000 * time.Millisecond)
		// Check if there are any neighbours
		currentNeighbours, currentConnectionsList := configuration.ReadCurrentResources("client.go: 29")
		if currentNeighbours == 0 {
			log.Println("No neighbours to send arbitrary transaction to")
			continue
		}
		// Log the current connections list
		log.Println("Current connections list:", currentConnectionsList.GetNodeConnections())
		// Send arbitrary transaction to all nodes
		propagation.SendArbitraryTransactionToAllNodes(currentConnectionsList)
	}
}
