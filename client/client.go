package client

import (
	"math/rand"
	"time"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/bootstrap"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/connection"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/nodestate"
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
		currentNeighbours, _ := nodestate.ReadCurrentConnections("client.go: 29")

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
	// selfNodeServerPort := configuration.GetSelfServerAddress().Port
	// if selfNodeServerPort != 0 {
	// 	go func() {
	// 		for {
	// 			// Wait for 10 seconds
	// 			time.Sleep(10000 * time.Millisecond)

	// 			// Display P2P Network
	// 			propagation.GetP2PNetwork()

	// 			// log.Println("Returned from DisplayP2PNetwork()")
	// 		}
	// 	}()
	// }

	// go CreateArbitraryTransactions()

	// Normal execution of the client
	for {
		// Wait for ten seconds. Then send arbitrary transaction to all nodes
		// // log.Println("Sleeping for three seconds")
		// time.Sleep(3000 * time.Millisecond)
		// // // log.Print("Awake")
		// // Check if there are any neighbours
		// currentNeighbours, _ := nodestate.ReadCurrentConnections("client.go: 66")
		// // // log.Println("Raed current connections")
		// if currentNeighbours == 0 {
		// 	// // log.Println("No neighbours to send arbitrary transaction to")
		// 	continue
		// }
		// // log the current connections list
		// // log.Println("Current connections list:", currentConnectionsList.GetNodeConnections())
		// Send arbitrary transaction to all nodes
		// propagation.SendArbitraryTransactionToAllNodes(currentConnectionsList)

		StartCreateBlocks()
	}
}
