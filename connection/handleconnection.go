package connection

import (
	"math/rand"
	"net"
	"time"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/bootstrap"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/nodestate"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

type NodesList = types.NodesList

func AddNewNodeConnection(ncl *ConnectionsList, nc NodeConnection, str string) bool {
	success := ncl.AddNodeConnection(nc)
	if success {
		// log.Printf("AddNewNodeConnection -> %s: Added new node connection with %v\n", str, nc.Node.GetAddress())
		go listenForMessages(nc)
	} else {
		// log.Printf("AddNewNodeConnection -> %s: Failed to add new node connection with %v\n", str, nc.Node.GetAddress())
		nc.Conn.Close()
	}
	return success
}

// Establishes connections with the nodes at the given addresses and returns the pointer to the connections' list
func establishConnectionWithNodes_SafeMode(existingNodesList NodesList) {
	minNeighbours := configuration.GetMinNeighbours()

	for _, node := range existingNodesList.GetNodes() {
		// // log.Println("Establishing connection with node", i+1, "at", node.GetAddress())
		currentNeighbours, currentConnections := nodestate.ReadCurrentConnections("handleconnection.go: 29")
		if currentNeighbours >= minNeighbours {
			// // log.Println("No need to add more connections. Currently established connections are", len(currentConnections.GetNodeConnections()), "neighbours")
			break
		}

		// Generate a random number between 0 and 1000
		randomTime := rand.Intn(1001)
		// Sleep for randomTime milliseconds
		time.Sleep(time.Duration(randomTime) * time.Millisecond)

		// Again check if the number of neighbours is greater than or equal to the minimum number of neighbours
		currentNeighbours, currentConnections = nodestate.ReadCurrentConnections("handleconnection.go: 29")
		if currentNeighbours >= minNeighbours {
			// // log.Println("No need to add more connections. Currently established connections are", len(currentConnections.GetNodeConnections()), "neighbours")
			break
		}

		// Check if this node is already connected
		alreadyConnected := currentConnections.ExistsAddress(node)

		if alreadyConnected {
			continue
		}

		// First establish a connection with the node
		var conn net.Conn = connectToNode(node)
		// // log.Println("Line 56: handleConnection.go")
		if conn == nil {
			// // log.Println("Line 58: handleConnection.go")
			// log.Println("ECWN_SM: Failed to establish connection with", node.GetAddress())
			continue
		}
		isConnectionRequestSuccess := sendConnectionRequestToNode(node, conn)
		// // log.Println("Line 59: handleConnection.go")
		if !isConnectionRequestSuccess {
			// log.Println("ECWN_SM: Failed to establish connection with", node.GetAddress())
			// Close the connection if the connection is not already closed
			if conn != nil {
				conn.Close()
			}
			continue
		}
		isConnectionResponseSuccess := receiveConnectionResponseFromNode(conn)
		if !isConnectionResponseSuccess {
			// log.Println("ECWN_SM: Failed to establish connection with", node.GetAddress())
			if conn != nil {
				conn.Close()
			}
			continue
		}
		nodeConnection := types.NewNodeConnection(node, conn)
		_, currentConnections = nodestate.LockCurrentConnections("Line 68: HandleCOnnecuin.go Add new node")
		returnVal := AddNewNodeConnection(&currentConnections, nodeConnection, "handleconnection")
		nodestate.UnlockCurrentConnections(currentConnections, "Line 70: UblockNode HandleCOnnecuin.go")
		if !returnVal {
			// log.Println("ECWN_SM: Failed to establish connection with", node.GetAddress())
			continue
		}
		// log.Println("ECWN_SM: Successfully established connection with", node.GetAddress())
		// // log.Println("Line 78: HandleConnection.go", node.GetAddress())
	}
}

func ConnectWithNetwork_SafeMode() {
	minNeighbours := configuration.GetMinNeighbours()
	// Read Current Resources
	currentNeighbours, _ := nodestate.ReadCurrentConnections("handleconnection.go: 76")
	if currentNeighbours >= minNeighbours {
		// log.Println("No need to add more connections. Currently established connections are", len(currentConnections.GetNodeConnections()), "neighbours")
		return
	}
	serverNode := configuration.GetSelfServerAddress()
	bootstrapNode := configuration.GetBootstrapNodeAddress()
	// Get the list of existing nodes in the network from the bootstrap node
	existingNodes := bootstrap.GetExistingNodesInNetwork(bootstrapNode, serverNode)
	if existingNodes == nil {
		// nodestate.LockBootstrapChan()
		// nodestate.UnlockBootstrapChan(false)
		// log.Println("Could not get the list of existing nodes in the network. Exiting...")
		return
	}
	existingNodesList := existingNodes.(NodesList)

	// // log.Println("Received existing nodes in network. Length:", len(existingNodes.GetNodes()))
	// // log.Println("Existing nodes in the network: ", len(existingNodesList.GetNodes()), existingNodesList.GetNodes())

	// Connect to the existing nodes
	establishConnectionWithNodes_SafeMode(existingNodesList)
}

func HandleLostNodeConnection(nc NodeConnection) {
	// Lock Resources
	_, currentConnections := nodestate.LockCurrentConnections("handleconnection.go: 99")
	currentConnections.RemoveNodeConnection(nc)
	// Unlock Resources
	nodestate.UnlockCurrentConnections(currentConnections, "handleconnection.go: 102")
	// if success {
	// 	// log.Println("Removed broken node connection with ", nc.Node.GetAddress())
	// } else {
	// 	// log.Println("Failed to remove broken node connection with ", nc.Node.GetAddress())
	// }
	ConnectWithNetwork_SafeMode()
}
