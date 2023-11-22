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
	// Register type string
	gob.Register("")
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
	// log.Println("Read current neighbours; value:", currentNeighbours)
	return currentNeighbours
}

// Non blocking read
func ReadCurrentConnections() types.ConnectionsList {
	currentConnections := <-currentConnectionsChan
	currentConnectionsChan <- currentConnections
	// log.Println("Read current connections; value:", currentConnections.GetNodeConnections())
	return currentConnections
}

// Critical Section
func LockCurrentNeighbours() uint32 {
	currentNeighbours := <-currentNeighboursChan
	// log.Println("Locking current neighbours; value:", currentNeighbours)
	return currentNeighbours
}

// Critical Section
func LockCurrentConnections() types.ConnectionsList {
	currentConnections := <-currentConnectionsChan
	// log.Println("Locking current connections; value:", currentConnections.GetNodeConnections())
	return currentConnections
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
	// log.Println("Unlocking current neighbours; value:", currentNeighbours)
	currentNeighboursChan <- currentNeighbours
}

// Critical Section
func UnlockCurrentConnections(currentConnections types.ConnectionsList) {
	// log.Println("Unlocking current connections; value:", currentConnections.GetNodeConnections())
	currentConnectionsChan <- currentConnections
}
