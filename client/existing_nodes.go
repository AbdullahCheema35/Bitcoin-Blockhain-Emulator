package client

import (
	"encoding/gob"
	"log"
	"net"
)

func getExistingNodesFromBootstrapNode(selfNode NodeAddress, conn net.Conn) *NodesList {
	// Send the node's server address to the bootstrap node
	enc := gob.NewEncoder(conn)
	err := enc.Encode(selfNode)
	if err != nil {
		log.Println("Error encoding:", err)
		return nil
	}

	// Receive the list of available nodes from the bootstrap node
	dec := gob.NewDecoder(conn)
	var existingNodes NodesList
	err = dec.Decode(&existingNodes)
	if err != nil {
		log.Println("Error decoding:", err)
		return nil
	}
	return &existingNodes
}

// Establishes a connection to the bootstrap server at the given address and returns the pointer to the connection
func connectToBootstrapNode(bootstrapNode NodeAddress) net.Conn {
	bootstrapAddress := bootstrapNode.GetAddress()
	conn, err := net.Dial("tcp", bootstrapAddress)
	if err != nil {
		log.Println("Error connecting to bootstrap node:", err)
		return nil
	}
	return conn
}

// returns the list of available nodes in the network from the bootstrap node
func getExistingNodesInNetwork(bootstrapNode NodeAddress, selfNode NodeAddress) *NodesList {
	// Connect to the bootstrap node
	var bootstrapConn net.Conn = connectToBootstrapNode(bootstrapNode)
	if bootstrapConn == nil {
		return nil
	}
	defer bootstrapConn.Close()

	// Get the list of available nodes from the bootstrap node
	var existingNodes *NodesList = getExistingNodesFromBootstrapNode(selfNode, bootstrapConn)
	return existingNodes
}
