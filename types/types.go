package types

import (
	"net"
	"strconv"
)

type NodeAddress struct {
	Port int
	IP   string
}

// NewNodeAddress creates a new NodeAddress instance with default IP value if not provided
// First parameter is port, second is NodeID (node's server addr), third is IP (optional)
func NewNodeAddress(port int, ip ...string) NodeAddress {
	var ipAddress string

	if len(ip) > 0 {
		ipAddress = ip[0]
	} else {
		ipAddress = "127.0.0.1" // Default IP if not provided
	}

	return NodeAddress{Port: port, IP: ipAddress}
}

// CombinedAddress returns the combined IP:Port as a string
func (na NodeAddress) GetAddress() string {
	return na.IP + ":" + strconv.Itoa(na.Port)
}

// NodesList represents a list of nodes
type NodesList struct {
	Nodes []NodeAddress
}

// NewNodesList creates a new NodesList instance
func NewNodesList() NodesList {
	return NodesList{Nodes: make([]NodeAddress, 0)}
}

// AddNode adds a node to the list
func (nl *NodesList) AddNode(node NodeAddress) {
	nl.Nodes = append(nl.Nodes, node)
}

// RemoveNode removes a node from the list
func (nl *NodesList) RemoveNode(node NodeAddress) {
	for i, n := range nl.Nodes {
		if n == node {
			nl.Nodes = append(nl.Nodes[:i], nl.Nodes[i+1:]...)
			break
		}
	}
}

// GetNodes returns the list of nodes
func (nl *NodesList) GetNodes() []NodeAddress {
	return nl.Nodes
}

type NodeConnection struct {
	Node NodeAddress
	Conn net.Conn
}

// NewNodeConnection creates a new NodeConnection instance
func NewNodeConnection(node NodeAddress, conn net.Conn) NodeConnection {
	return NodeConnection{Node: node, Conn: conn}
}

// ConnectionsList represents a list of node connections
type ConnectionsList struct {
	NodeConnections []NodeConnection
}

// NewNodeConnectionsList creates a new ConnectionsList instance
func NewNodeConnectionsList() ConnectionsList {
	return ConnectionsList{NodeConnections: make([]NodeConnection, 0)}
}

// AddNodeConnection adds a node connection to the list
func (ncl *ConnectionsList) AddNodeConnection(nodeConnection NodeConnection) {
	ncl.NodeConnections = append(ncl.NodeConnections, nodeConnection)
}

// RemoveNodeConnection removes a node connection from the list
func (ncl *ConnectionsList) RemoveNodeConnection(nodeConnection NodeConnection) {
	for i, n := range ncl.NodeConnections {
		if n == nodeConnection {
			ncl.NodeConnections = append(ncl.NodeConnections[:i], ncl.NodeConnections[i+1:]...)
			break
		}
	}
}

// GetNodeConnections returns the list of node connections
func (ncl *ConnectionsList) GetNodeConnections() []NodeConnection {
	return ncl.NodeConnections
}

type ConnectionRequestType uint8

const (
	ConnectionRequestTypeSuccess ConnectionRequestType = iota
	ConnectionRequestTypeFailure
)

type ConnectionResponseType uint8

const (
	ConnectionResponseTypeSuccess ConnectionResponseType = iota
	ConnectionResponseTypeFailure
)

type MessageType uint8

type MessageHeader struct {
	Type   MessageType
	Sender NodeAddress
}

const (
	MessageTypeTransaction MessageType = iota
	MessageTypeBlock
	MessageTypeRequest
	MessageTypeResponse
	MessageTypeConnection
	MessageTypeConnectionRequest
	MessageTypeConnectionResponse
)

type Message struct {
	Header MessageHeader
	Body   interface{}
	Sender interface{}
}

func NewMessage(header MessageHeader, body interface{}, sender interface{}) Message {
	return Message{Header: header, Body: body, Sender: sender}
}
