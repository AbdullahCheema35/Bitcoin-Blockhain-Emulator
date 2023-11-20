package configuration

import "github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"

var (
	minNeighbours uint32
	maxNeighbours uint32
)

var (
	currentNeighboursChan  chan uint32                = make(chan uint32, 1)
	currentConnectionsChan chan types.ConnectionsList = make(chan types.ConnectionsList, 1)
)

var (
	selfServerAddress    types.NodeAddress
	selfBootstrapAddress types.NodeAddress
	bootstrapNodeAddress types.NodeAddress
	isSelfBootstrapNode  bool
)

func InitConfiguration(_minNeighbours, _maxNeighbours uint32, _selfServerAddress, _selfBootstrapAddress, _bootstrapNodeAddress types.NodeAddress, _isSelfBootstrapNode bool) {
	minNeighbours = _minNeighbours
	maxNeighbours = _maxNeighbours

	currentNeighboursChan <- 0
	currentConnectionsChan <- types.NewNodeConnectionsList()

	selfServerAddress = _selfServerAddress
	selfBootstrapAddress = _selfBootstrapAddress
	bootstrapNodeAddress = _bootstrapNodeAddress
	isSelfBootstrapNode = _isSelfBootstrapNode
}

// Getter functions
func GetMinNeighbours() uint32 {
	return minNeighbours
}

func GetMaxNeighbours() uint32 {
	return maxNeighbours
}

func GetCurrentNeighbours() uint32 {
	return <-currentNeighboursChan
}

func GetCurrentConnections() types.ConnectionsList {
	return <-currentConnectionsChan
}

func GetSelfServerAddress() types.NodeAddress {
	return selfServerAddress
}

func GetSelfBootstrapAddress() types.NodeAddress {
	return selfBootstrapAddress
}

func GetBootstrapNodeAddress() types.NodeAddress {
	return bootstrapNodeAddress
}

func GetIsSelfBootstrapNode() bool {
	return isSelfBootstrapNode
}

// Setter functions for channels
func SetCurrentNeighbours(currentNeighbours uint32) {
	currentNeighboursChan <- currentNeighbours
}

func SetCurrentConnections(currentConnections types.ConnectionsList) {
	currentConnectionsChan <- currentConnections
}
