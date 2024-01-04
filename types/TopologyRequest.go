package types

import "fmt"

// NetworkList represents a list of network connections between nodes
type NetworkList struct {
	Src NodeAddress
	Dst NodesList
}

// NewNetworkList creates a new NetworkList
func NewNetworkList(src NodeAddress, dst NodesList) NetworkList {
	return NetworkList{
		Src: src,
		Dst: dst,
	}
}

type TopologyRequest struct {
	Origin        NodesList
	NodesFound    NodesList
	ThisNodePeers NetworkList
}

func NewTopologyRequest(origin NodesList, nodesFound NodesList, thisNodePeers NetworkList) TopologyRequest {
	return TopologyRequest{
		Origin:        origin,
		NodesFound:    nodesFound,
		ThisNodePeers: thisNodePeers,
	}
}

// Display prints a text-based UI showing the connections of every node in the network list
func (nl NetworkList) Display() {
	// Display the connections
	fmt.Printf("Node %s\n", nl.Src.GetAddress())
	fmt.Println("+- Connections:")

	for _, node := range nl.Dst.GetNodes() {
		if node != nl.Src {
			fmt.Printf("   +- %s\n", node.GetAddress())
		}
	}
}
