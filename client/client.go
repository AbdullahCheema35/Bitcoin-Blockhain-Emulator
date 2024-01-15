package client

import (
	"math/rand"
	"time"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/bootstrap"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/connection"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/mineblock"
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

	// Start creating blocks
	var miningChan chan bool
	for {
		if EnoughTransactionsForBlock() {
			miningChan = nodestate.InitMiningChan()
			ok, newCreatedBlock := createBlock()
			if ok {
				mineblock.MineBlock(newCreatedBlock, miningChan)
			}
		}
		time.Sleep(1000 * time.Millisecond)
	}
}
