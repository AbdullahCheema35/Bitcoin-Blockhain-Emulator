package types

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
