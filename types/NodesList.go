package types

// NodesList represents a list of nodes
type NodesList struct {
	Nodes []NodeAddress
}

// NewNodesList creates a new NodesList instance
func NewNodesList() NodesList {
	return NodesList{Nodes: make([]NodeAddress, 0)}
}

// AddNode adds a node to the list
func (nl *NodesList) AddNode(node NodeAddress) bool {
	for _, n := range nl.Nodes {
		if n.GetAddress() == node.GetAddress() {
			return false
		}
	}
	nl.Nodes = append(nl.Nodes, node)
	return true
}

// RemoveNode removes a node from the list
func (nl *NodesList) RemoveNode(node NodeAddress) bool {
	for i, n := range nl.Nodes {
		if n == node {
			nl.Nodes = append(nl.Nodes[:i], nl.Nodes[i+1:]...)
			return true
		}
	}
	return false
}

// GetNodes returns the list of nodes
func (nl *NodesList) GetNodes() []NodeAddress {
	return nl.Nodes
}
