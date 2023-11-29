package types

import (
	"fmt"
	"time"
)

type BootstrapNodeAddress struct {
	Node        NodeAddress
	LastUpdated time.Time
}

// NewBootstrapNodeAddress creates a new BootstrapNodeAddress instance
func NewBootstrapNodeAddress(node NodeAddress, lastUpdated time.Time) BootstrapNodeAddress {
	return BootstrapNodeAddress{Node: node, LastUpdated: lastUpdated}
}

func PrintBootstrapNodeAddress(bootstrapNode BootstrapNodeAddress) {
	fmt.Println("Bootstrap Node Address:")
	fmt.Printf("Node: %s\n", bootstrapNode.Node.GetAddress())
	fmt.Printf("Last Updated: %s\n", bootstrapNode.LastUpdated.Format(time.RFC3339))
	fmt.Println("------------------------")
}
