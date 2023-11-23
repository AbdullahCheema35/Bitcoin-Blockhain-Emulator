package types

import "time"

type BootstrapNodeAddress struct {
	Node        NodeAddress
	LastUpdated time.Time
}

// NewBootstrapNodeAddress creates a new BootstrapNodeAddress instance
func NewBootstrapNodeAddress(node NodeAddress, lastUpdated time.Time) BootstrapNodeAddress {
	return BootstrapNodeAddress{Node: node, LastUpdated: lastUpdated}
}
