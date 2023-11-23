package configuration

import (
	"encoding/gob"
	"math/rand"
	"time"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

const (
	minNeighbours       int = 2
	maxNeighbours       int = 4
	maxSecondsPingDelay int = 12
)

var (
	currentConnectionsChan   chan types.ConnectionsList   = make(chan types.ConnectionsList, 1)
	currentExistingNodesChan chan types.BootstrapNodesMap = make(chan types.BootstrapNodesMap, 1)
	bootstrapLockChan        chan bool                    = make(chan bool, 1)
)

var (
	selfServerAddress    types.NodeAddress
	selfBootstrapAddress types.NodeAddress
	bootstrapNodeAddress types.NodeAddress
	isSelfBootstrapNode  bool
	// connectedWithBootstrap bool
)

func InitConfiguration(_selfServerAddress, _selfBootstrapAddress, _bootstrapNodeAddress types.NodeAddress, _isSelfBootstrapNode bool) {
	newConnectionsList := types.NewNodeConnectionsList()
	newExistingNodesMap := types.NewBootstrapNodesMap()
	currentConnectionsChan <- newConnectionsList
	currentExistingNodesChan <- newExistingNodesMap
	// connectedWithBootstrap = false
	bootstrapLockChan <- false

	selfServerAddress = _selfServerAddress
	selfBootstrapAddress = _selfBootstrapAddress
	bootstrapNodeAddress = _bootstrapNodeAddress
	isSelfBootstrapNode = _isSelfBootstrapNode

	// Register types for gob
	gob.Register(types.ConnectionRequestTypeFailure)
	gob.Register(types.ConnectionResponseTypeSuccess)
	gob.Register(types.NodesList{})
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

func GetMaxSecondsPingDelay() int {
	return maxSecondsPingDelay
}

func LockBootstrapChan() bool {
	return <-bootstrapLockChan
}

func UnlockBootstrapChan(connected bool) {
	bootstrapLockChan <- connected
}

// Reader functions for channels
// Non blocking read
func ReadCurrentConnections(line string) (int, types.ConnectionsList) {
	currentNeighbours, currentConnections := LockCurrentConnections("configuration.go: 56")
	UnlockCurrentConnections(currentConnections, "configuration.go: 57")
	return currentNeighbours, currentConnections
}

func ReadCurrentExistingNodes(line string) (int, types.BootstrapNodesMap) {
	currentNeighbours, currentExistingNodes := LockCurrentExistingNodes("configuration.go: 62")
	UnlockCurrentExistingNodes(currentExistingNodes, "configuration.go: 63")
	return currentNeighbours, currentExistingNodes
}

// Reader functions for channels
// Blocking read
func LockCurrentConnections(line string) (int, types.ConnectionsList) {
	// log.Println("Locking current resources", line)
	currentConnections := <-currentConnectionsChan
	currentNeighbours := len(currentConnections.GetNodeConnections())
	// // log.Println("Locking current neighbours; value:", currentNeighbours)
	// // log.Println("Locking current connections; value:", currentConnections.GetNodeConnections())
	return currentNeighbours, currentConnections
}

func LockCurrentExistingNodes(line string) (int, types.BootstrapNodesMap) {
	// log.Println("Locking current resources", line)
	currentExistingNodes := <-currentExistingNodesChan
	currentNeighbours := currentExistingNodes.GetLength()
	// // log.Println("Locking current neighbours; value:", currentNeighbours)
	// // log.Println("Locking current connections; value:", currentConnections.GetNodeConnections())
	return currentNeighbours, currentExistingNodes
}

// Writer functions for channels
// Blocking write (if the channel is full, i.e., above the buffer size)
func UnlockCurrentConnections(currentConnections types.ConnectionsList, line string) {
	// // log.Println("Unlocking current neighbours; value:", currentNeighbours)
	// // log.Println("Unlocking current connections; value:", currentConnections.GetNodeConnections())
	currentConnectionsChan <- currentConnections
	// log.Println("Unlocked current resources", line)
}

func UnlockCurrentExistingNodes(currentExistingNodes types.BootstrapNodesMap, line string) {
	// // log.Println("Unlocking current neighbours; value:", currentNeighbours)
	// // log.Println("Unlocking current connections; value:", currentConnections.GetNodeConnections())
	currentExistingNodesChan <- currentExistingNodes
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
