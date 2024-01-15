package types

import "strconv"

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
