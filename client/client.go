package client

import (
	"log"
	"math/rand"
	"time"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/bootstrap"
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

func handleBootstrapServer() {
	bootstrapNode := configuration.GetBootstrapNodeAddress()

	// Generate a random number between 0 and 100
	minTime, maxTime := 100, 200
	randomTime := minTime + rand.Intn((maxTime-minTime)+1)
	// Sleep for randomTime milliseconds
	time.Sleep(time.Duration(randomTime) * time.Millisecond)

	for {
		// Read the current connections list
		currentNeighbours, _ := configuration.ReadCurrentConnections("client.go: 29")

		// If the number of neighbours is less than the minimum number of neighbours, then connect to the network
		if currentNeighbours < configuration.GetMinNeighbours() {
			connection.ConnectWithNetwork_SafeMode()
		}

		// Go to sleep for 5 seconds
		time.Sleep(5000 * time.Millisecond)

		// Ping the bootstrap server
		bootstrap.PingBootstrapServer(bootstrapNode)
	}
}

func StartClient() {
	// Start handling bootstrap server node
	go handleBootstrapServer()

	// Normal execution of the client
	for {
		// Wait for ten seconds. Then send arbitrary transaction to all nodes
		// log.Println("Waiting for one seconds")
		time.Sleep(30000 * time.Millisecond)
		// Check if there are any neighbours
		currentNeighbours, currentConnectionsList := configuration.ReadCurrentConnections("client.go: 29")
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
