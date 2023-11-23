package types

import "time"

// BootstrapNodesMap represents a list of bootstrap nodes
type BootstrapNodesMap struct {
	BootstrapNodes map[string]BootstrapNodeAddress
}

// NewBootstrapNodesMap creates a new BootstrapNodesMap instance
func NewBootstrapNodesMap() BootstrapNodesMap {
	return BootstrapNodesMap{BootstrapNodes: make(map[string]BootstrapNodeAddress)}
}

// AddBootstrapNode adds a bootstrap node to the list
func (bnl *BootstrapNodesMap) AddUpdateBootstrapNode(nodeAddr NodeAddress) {
	bnl.BootstrapNodes[nodeAddr.GetAddress()] = NewBootstrapNodeAddress(nodeAddr, time.Now())
}

// RemoveBootstrapNode removes a bootstrap node from the map
func (bnl *BootstrapNodesMap) RemoveBootstrapNode(nodeAddr NodeAddress) {
	delete(bnl.BootstrapNodes, nodeAddr.GetAddress())
}

// // GetBootstrapNodes returns a list of bootstrap nodes
// func (bnl *BootstrapNodesMap) GetBootstrapNodesList() []BootstrapNodeAddress {
// 	bootstrapNodes := make([]BootstrapNodeAddress, 0)
// 	for _, node := range bnl.BootstrapNodes {
// 		bootstrapNodes = append(bootstrapNodes, node)
// 	}
// 	return bootstrapNodes
// }

// GetRecentBootstrapNodes returns a list of bootstrap nodes that were updated recently
func (bnl *BootstrapNodesMap) GetRecentBootstrapNodes(maxSecondsPingDelay int) NodesList {
	recentNodesList := NewNodesList()
	for _, node := range bnl.BootstrapNodes {
		if time.Since(node.LastUpdated).Seconds() <= float64(maxSecondsPingDelay) {
			recentNodesList.AddNode(node.Node)
		} else {
			// Remove this node from the map
			bnl.RemoveBootstrapNode(node.Node)
		}
	}
	return recentNodesList
}

// Get length of bootstrap nodes map
func (bnl *BootstrapNodesMap) GetLength() int {
	return len(bnl.BootstrapNodes)
}
