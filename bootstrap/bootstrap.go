package bootstrap

import (
	"encoding/gob"
	"log"
	"net"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

type NodeAddress = types.NodeAddress

type NodesList = types.NodesList

// Handle the list of available nodes
// nListChan is the write-only channel to send the list of available nodes to handleBootstrapQuery goroutine
// nAddrChan is the read-only channel to receive the querying node's server address from handleBootstrapQuery goroutine
func handleAvailableNodesList(nListChan chan<- NodesList, nAddrChan <-chan NodeAddress, sNode NodeAddress) {
	var availableNodes NodesList = types.NewNodesList()

	// Add the bootstrap node's normal server address to the list of available nodes
	availableNodes.AddNode(sNode)

	// Loop until the channel is closed
	for nodeAddress := range nAddrChan {
		// Send the list of available nodes to the goroutine that handles the query from a node's server to the bootstrap node
		nListChan <- availableNodes

		// Add the new node's server address to the list of available nodes
		availableNodes.AddNode(nodeAddress)

		log.Println("Added", nodeAddress.GetAddress(), "to the list of available nodes")
	}
}

// Handle the query from a node's server to the bootstrap node
// conn is the connection to the querying node's server
// nListChan is the read-only channel to receive the list of available nodes from handleAvailableNodesList goroutine
// nAddrChan is the write-only channel to send the querying node's server address to handleAvailableNodesList goroutine
func handleBootstrapQuery(conn net.Conn, nListChan <-chan NodesList, nAddrChan chan<- NodeAddress) {
	defer conn.Close()

	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(conn)

	// Receive the NodeAddress of the node's server that is querying to the bootstrap node
	var nodeAddress NodeAddress
	err := dec.Decode(&nodeAddress)
	if err != nil {
		log.Println("Error decoding:", err)
		return
	}
	log.Println("Received bootstrap query from", nodeAddress.GetAddress())

	// First, send the querying node's server address to the goroutine that handles the list of available nodes
	nAddrChan <- nodeAddress

	// Then, receive the list of available nodes from the goroutine that handles the list of available nodes
	availableNodes := <-nListChan

	// Send the list of available nodes to the querying node
	err = enc.Encode(availableNodes)
	if err != nil {
		log.Println("Error encoding:", err)
		return
	}

	log.Println("Sent list of available nodes to", nodeAddress.GetAddress())
}

func StartBootstrapServer(bNode NodeAddress, sNode NodeAddress) {
	bootstrapAddress := bNode.GetAddress()
	listener, err := net.Listen("tcp", bootstrapAddress)
	if err != nil {
		log.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	log.Println("Bootstrap node listening on port", bNode.Port)

	// Create channel for NodesList to be shared between goroutines
	nodesListChannel := make(chan NodesList)
	defer close(nodesListChannel)

	// Create channel for NodeAddress to be shared between goroutines
	nodeAddressChannel := make(chan NodeAddress)
	defer close(nodeAddressChannel)

	// Start goroutine to handle the list of available nodes
	go handleAvailableNodesList(nodesListChannel, nodeAddressChannel, sNode)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go handleBootstrapQuery(conn, nodesListChannel, nodeAddressChannel)
	}
}
