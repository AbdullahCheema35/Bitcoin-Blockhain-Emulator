package bootstrap

import (
	"encoding/gob"
	"log"
	"net"
	"time"
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

	// filter the list of available nodes to remove the node's own server address
	existingNodes.RemoveNode(selfNode)

	return &existingNodes
}

// Establishes a connection to the bootstrap server at the given address and returns the pointer to the connection
func connectToBootstrapNode(bootstrapNode NodeAddress) net.Conn {
	bootstrapAddress := bootstrapNode.GetAddress()
	for seconds := 1; seconds < 100; seconds *= 2 {
		conn, err := net.DialTimeout("tcp", bootstrapAddress, 1000*time.Millisecond)
		if err != nil {
			// log.Println("Couldn't connect to bootstrap node:", err, "Retrying in ", seconds, "seconds")
			// Wait for a while before retrying
			time.Sleep(time.Duration(seconds) * time.Second)
		} else {
			// log.Println("Connected to bootstrap node")
			return conn
		}
	}
	return nil
}

// returns the list of available nodes in the network from the bootstrap node
func GetExistingNodesInNetwork(bootstrapNode NodeAddress, selfNode NodeAddress) *NodesList {
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
