package configuration

import (
	"encoding/gob"
	"math/rand"
	"time"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

const (
	minNeighbours int = 2
	maxNeighbours int = 4
)

var (
	currentNeighboursChan  chan int                   = make(chan int, 1)
	currentConnectionsChan chan types.ConnectionsList = make(chan types.ConnectionsList, 1)
)

var (
	selfServerAddress    types.NodeAddress
	selfBootstrapAddress types.NodeAddress
	bootstrapNodeAddress types.NodeAddress
	isSelfBootstrapNode  bool
)

func InitConfiguration(_selfServerAddress, _selfBootstrapAddress, _bootstrapNodeAddress types.NodeAddress, _isSelfBootstrapNode bool) {
	newConnectionsList := types.NewNodeConnectionsList()
	currentNeighboursChan <- len(newConnectionsList.GetNodeConnections())
	currentConnectionsChan <- newConnectionsList

	selfServerAddress = _selfServerAddress
	selfBootstrapAddress = _selfBootstrapAddress
	bootstrapNodeAddress = _bootstrapNodeAddress
	isSelfBootstrapNode = _isSelfBootstrapNode

	// Register types for gob
	gob.Register(types.ConnectionRequestTypeFailure)
	gob.Register(types.ConnectionResponseTypeSuccess)
	// Register type string
	gob.Register("")

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())
}

// Getter functions
func GetMinNeighbours() int {
	return minNeighbours
}

func GetMaxNeighbours() int {
	return maxNeighbours
}

// Reader functions for channels
// Non blocking read
func ReadCurrentResources(line string) (int, types.ConnectionsList) {
	currentNeighbours, currentConnections := LockCurrentResources("configuration.go: 56")
	UnlockCurrentResources(currentConnections, "configuration.go: 57")
	return currentNeighbours, currentConnections
}

func LockCurrentResources(line string) (int, types.ConnectionsList) {
	// log.Println("Locking current resources", line)
	currentNeighbours := <-currentNeighboursChan
	currentConnections := <-currentConnectionsChan
	// // log.Println("Locking current neighbours; value:", currentNeighbours)
	// // log.Println("Locking current connections; value:", currentConnections.GetNodeConnections())
	return currentNeighbours, currentConnections
}

func UnlockCurrentResources(currentConnections types.ConnectionsList, line string) {
	currentNeighbours := len(currentConnections.GetNodeConnections())
	// // log.Println("Unlocking current neighbours; value:", currentNeighbours)
	// // log.Println("Unlocking current connections; value:", currentConnections.GetNodeConnections())
	// typecast to int
	currentNeighboursChan <- currentNeighbours
	currentConnectionsChan <- currentConnections
	// log.Println("Unlocked current resources", line)
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
