package configuration

import (
	"encoding/gob"
	"math/rand"
	"time"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

var (
	minNeighbours          int = 2
	maxNeighbours          int = 4
	maxSecondsPingDelay    int = 12
	maxTransactionsInBlock int = 5
	minTransactionsInBlock int = 2
	difficultyTarget       int = 24
	transactionSpeed       int = 10
)

var (
	selfServerAddress    types.NodeAddress
	selfBootstrapAddress types.NodeAddress
	bootstrapNodeAddress types.NodeAddress
	isSelfBootstrapNode  bool
)

func InitConfiguration(_selfServerAddress, _selfBootstrapAddress, _bootstrapNodeAddress types.NodeAddress, _isSelfBootstrapNode bool) {
	selfServerAddress = _selfServerAddress
	selfBootstrapAddress = _selfBootstrapAddress
	bootstrapNodeAddress = _bootstrapNodeAddress
	isSelfBootstrapNode = _isSelfBootstrapNode

	// Register types for gob
	gob.Register(types.NodesList{})
	// Register type string
	gob.Register("")
	// Register type Transaction
	gob.Register(types.Transaction{})
	// Register type Block
	gob.Register(types.Block{})
	// Register type BlockChain
	gob.Register(types.BlockChain{})
	// Register type TopologyRequest
	gob.Register(types.TopologyRequest{})

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

func GetMaxTransactionsInBlock() int {
	return maxTransactionsInBlock
}

func GetMinTransactionsInBlock() int {
	return minTransactionsInBlock
}

func GetTransactionSpeed() int {
	return transactionSpeed
}

func SetTransactionSpeed(value int) { transactionSpeed = value }

func GetDifficultyTarget() int {
	return difficultyTarget
}

func SetMinNeighbours(value int) { minNeighbours = value }

func SetMaxNeighbours(value int) { maxNeighbours = value }

func SetMaxSecondsPingDelay(value int) { maxSecondsPingDelay = value }

func SetMaxTransactionsInBlock(value int) { maxTransactionsInBlock = value }

func SetMinTransactionsInBlock(value int) { minTransactionsInBlock = value }

func SetDifficultyTarget(value int) { difficultyTarget = value }
