package connection

import (
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/bootstrap"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

type NodesList = types.NodesList

func AddNewNodeConnection(ncl *ConnectionsList, nc NodeConnection, str string) bool {
	success := ncl.AddNodeConnection(nc)
	if success {
		log.Printf("AddNewNodeConnection -> %s: Added new node connection with %v\n", str, nc.Node.GetAddress())
		go listenForMessages(nc)
	} else {
		log.Printf("AddNewNodeConnection -> %s: Failed to add new node connection with %v\n", str, nc.Node.GetAddress())
		nc.Conn.Close()
	}
	return success
}

// Establishes connections with the nodes at the given addresses and returns the pointer to the connections' list
func establishConnectionWithNodes_SafeMode(existingNodesList *NodesList) {
	minNeighbours := configuration.GetMinNeighbours()

	for _, node := range existingNodesList.GetNodes() {
		// log.Println("Establishing connection with node", i+1, "at", node.GetAddress())
		// Generate a random number between 0 and 1000
		randomTime := rand.Intn(1001)
		// Sleep for randomTime milliseconds
		time.Sleep(time.Duration(randomTime) * time.Millisecond)

		currentNeighbours, currentConnections := configuration.ReadCurrentResources("handleconnection.go: 29")
		if currentNeighbours >= minNeighbours {
			// log.Println("No need to add more connections. Currently established connections are", len(currentConnections.GetNodeConnections()), "neighbours")
			break
		}

		alreadyConnected := currentConnections.ExistsAddress(node)

		if alreadyConnected {
			continue
		}

		successInProcess := false
		// First establish a connection with the node
		var conn net.Conn = connectToNode(node)
		// log.Println("Line 56: handleConnection.go")
		if conn != nil {
			// Send a connection request to the node
			isConnectionRequestSuccess := sendConnectionRequestToNode(node, conn)
			// log.Println("Line 60: handleConnection.go")
			if isConnectionRequestSuccess {
				// Receive a connection response from the node
				// log.Println("Line 63: handleConnection.go")
				isConnectionResponseSuccess := receiveConnectionResponseFromNode(conn)
				// log.Println("Line 65: handleConnection.go")
				if isConnectionResponseSuccess { // Successfully established connection with the node
					// Add the node to the list of current neighbours
					// log.Println("Line 68: HandleConnection.go", node.GetAddress())
					nodeConnection := types.NewNodeConnection(node, conn)
					// log.Println("Line 70: HandleConnection.go", node.GetAddress())
					_, currentConnections := configuration.LockCurrentResources("Line 68: HandleCOnnecuin.go Add new node")
					returnVal := AddNewNodeConnection(&currentConnections, nodeConnection, "handleconnection")
					configuration.UnlockCurrentResources(currentConnections, "Line 70: UblockNode HandleCOnnecuin.go")

					if returnVal {
						successInProcess = true
					}
					// log.Println("Line 72: HandleConnection.go", node.GetAddress())
				}
				// log.Println("Line 74: HandleConnection.go", node.GetAddress())
			}
			// log.Println("Line 76: HandleConnection.go", node.GetAddress())
		}
		if successInProcess {
			log.Println("ECWN_SM: Successfully established connection with", node.GetAddress())
		} else {
			log.Println("ECWN_SM: Failed to establish connection with", node.GetAddress())
		}
		// log.Println("Line 78: HandleConnection.go", node.GetAddress())
	}
}

func ConnectWithNetwork_SafeMode() {
	minNeighbours := configuration.GetMinNeighbours()
	// Read Current Resources
	currentNeighbours, currentConnections := configuration.ReadCurrentResources("handleconnection.go: 76")
	if currentNeighbours >= minNeighbours {
		log.Println("No need to add more connections. Currently established connections are", len(currentConnections.GetNodeConnections()), "neighbours")
		return
	}
	serverNode := configuration.GetSelfServerAddress()
	bootstrapNode := configuration.GetBootstrapNodeAddress()
	// Get the list of existing nodes in the network from the bootstrap node
	var existingNodes *NodesList = bootstrap.GetExistingNodesInNetwork(bootstrapNode, serverNode)
	if existingNodes == nil {
		log.Println("Could not get the list of existing nodes in the network. Exiting Client...")
		return
	}

	// log.Println("Received existing nodes in network. Length:", len(existingNodes.GetNodes()))
	log.Println("Existing nodes in the network: ", len(existingNodes.GetNodes()), existingNodes.GetNodes())

	// Connect to the existing nodes
	establishConnectionWithNodes_SafeMode(existingNodes)
}

func HandleLostNodeConnection(nc NodeConnection) {
	// Lock Resources
	_, currentConnections := configuration.LockCurrentResources("handleconnection.go: 99")
	success := currentConnections.RemoveNodeConnection(nc)
	// Unlock Resources
	configuration.UnlockCurrentResources(currentConnections, "handleconnection.go: 102")
	if success {
		log.Println("Removed broken node connection with ", nc.Node.GetAddress())
	} else {
		log.Println("Failed to remove broken node connection with ", nc.Node.GetAddress())
	}
	ConnectWithNetwork_SafeMode()
}
