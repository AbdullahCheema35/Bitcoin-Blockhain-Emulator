package types

// ConnectionsList represents a list of node connections
type ConnectionsList struct {
	NodeConnections []NodeConnection
}

// NewNodeConnectionsList creates a new ConnectionsList instance
func NewNodeConnectionsList() ConnectionsList {
	return ConnectionsList{NodeConnections: make([]NodeConnection, 0)}
}

// AddNodeConnection adds a node connection to the list
func (ncl *ConnectionsList) AddNodeConnection(nodeConnection NodeConnection) bool {
	for _, n := range ncl.NodeConnections {
		if n.Node == nodeConnection.Node {
			return false
		}
	}
	ncl.NodeConnections = append(ncl.NodeConnections, nodeConnection)
	return true
}

// RemoveNodeConnection removes a node connection from the list
func (ncl *ConnectionsList) RemoveNodeConnection(nodeConnection NodeConnection) bool {
	for i, n := range ncl.NodeConnections {
		if n.Node == nodeConnection.Node {
			ncl.NodeConnections = append(ncl.NodeConnections[:i], ncl.NodeConnections[i+1:]...)
			return true
		}
	}
	return false
}

func (ncl *ConnectionsList) ExistsAddress(nodeAddress NodeAddress) bool {
	for _, n := range ncl.NodeConnections {
		if n.Node == nodeAddress {
			return true
		}
	}
	return false
}

// GetNodeConnections returns the list of node connections
func (ncl *ConnectionsList) GetNodeConnections() []NodeConnection {
	return ncl.NodeConnections
}
