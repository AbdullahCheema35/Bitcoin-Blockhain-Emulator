package configuration

import (
	"encoding/gob"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

const (
	minNeighbours uint32 = 2
	maxNeighbours uint32 = 4
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

func InitConfiguration(_selfServerAddress, _selfBootstrapAddress, _bootstrapNodeAddress types.NodeAddress, _isSelfBootstrapNode bool) {
	currentNeighboursChan <- 0
	currentConnectionsChan <- types.NewNodeConnectionsList()

	selfServerAddress = _selfServerAddress
	selfBootstrapAddress = _selfBootstrapAddress
	bootstrapNodeAddress = _bootstrapNodeAddress
	isSelfBootstrapNode = _isSelfBootstrapNode

	// Register types for gob
	gob.Register(types.ConnectionRequestTypeFailure)
	gob.Register(types.ConnectionResponseTypeSuccess)
	gob.Register(types.MessageTypeTransaction)
}

// Getter functions
func GetMinNeighbours() uint32 {
	return minNeighbours
}

func GetMaxNeighbours() uint32 {
	return maxNeighbours
}

// Reader functions for channels
// Non blocking read
func ReadCurrentNeighbours() uint32 {
	currentNeighbours := <-currentNeighboursChan
	currentNeighboursChan <- currentNeighbours
	return currentNeighbours
}

// Non blocking read
func ReadCurrentConnections() types.ConnectionsList {
	currentConnections := <-currentConnectionsChan
	currentConnectionsChan <- currentConnections
	return currentConnections
}

// Critical Section
func LockCurrentNeighbours() uint32 {
	return <-currentNeighboursChan
}

// Critical Section
func LockCurrentConnections() types.ConnectionsList {
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
// Critical Section
func UnlockCurrentNeighbours(currentNeighbours uint32) {
	currentNeighboursChan <- currentNeighbours
}

// Critical Section
func UnlockCurrentConnections(currentConnections types.ConnectionsList) {
	currentConnectionsChan <- currentConnections
}
