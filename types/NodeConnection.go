package types

import "net"

// NodeConnection represents a connection with a node
type NodeConnection struct {
	Node NodeAddress
	Conn net.Conn
}

// NewNodeConnection creates a new NodeConnection instance
func NewNodeConnection(node NodeAddress, conn net.Conn) NodeConnection {
	return NodeConnection{Node: node, Conn: conn}
}
