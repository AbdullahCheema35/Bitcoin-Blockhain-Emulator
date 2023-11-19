package client

import (
	"log"
	"net"

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

// Establishes connections with the nodes at the given addresses and returns the pointer to the connections' list
func establishConnectionWithExistingNodes(existingNodesList NodesList) ConnectionsList {
	var nodeConnectionsList ConnectionsList = types.NewNodeConnectionsList()
	minNeighbours := configuration.MinNeighbours
	currentNeighbours := configuration.CurrentNeighbours

	for _, node := range existingNodesList.GetNodes() {
		if currentNeighbours >= minNeighbours {
			break
		}
		var conn net.Conn = connectToNode(node)
		if conn != nil {
			currentNeighbours++
			nodeConnection := types.NewNodeConnection(node, conn)
			nodeConnectionsList.AddNodeConnection(nodeConnection)
		}
	}
	return nodeConnectionsList
}
