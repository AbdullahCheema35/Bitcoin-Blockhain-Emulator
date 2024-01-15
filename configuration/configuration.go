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
	minTransactionsInBlock int = 3
	difficultyTarget       int = 4
	transactionSpeed       int = 5
)

// Create individual channels for all of above variables
var (
	minNeighboursChan          = make(chan bool, 1)
	maxNeighboursChan          = make(chan bool, 1)
	maxSecondsPingDelayChan    = make(chan bool, 1)
	maxTransactionsInBlockChan = make(chan bool, 1)
	minTransactionsInBlockChan = make(chan bool, 1)
	difficultyTargetChan       = make(chan bool, 1)
	transactionSpeedChan       = make(chan bool, 1)
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

	// Initialize the channels
	minNeighboursChan <- true
	maxNeighboursChan <- true
	maxSecondsPingDelayChan <- true
	maxTransactionsInBlockChan <- true
	minTransactionsInBlockChan <- true
	difficultyTargetChan <- true
	transactionSpeedChan <- true

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
	<-minNeighboursChan
	retVal := minNeighbours
	minNeighboursChan <- true
	return retVal
}

func GetMaxNeighbours() int {
	<-maxNeighboursChan
	retVal := maxNeighbours
	maxNeighboursChan <- true
	return retVal
}

func GetMaxSecondsPingDelay() int {
	<-maxSecondsPingDelayChan
	retVal := maxSecondsPingDelay
	maxSecondsPingDelayChan <- true
	return retVal
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
	<-maxTransactionsInBlockChan
	retVal := maxTransactionsInBlock
	maxTransactionsInBlockChan <- true
	return retVal
}

func GetMinTransactionsInBlock() int {
	<-minTransactionsInBlockChan
	retVal := minTransactionsInBlock
	minTransactionsInBlockChan <- true
	return retVal
}

func GetTransactionSpeed() int {
	<-transactionSpeedChan
	retVal := transactionSpeed
	transactionSpeedChan <- true
	return retVal
}

func SetTransactionSpeed(value int) {
	<-transactionSpeedChan
	transactionSpeed = value
	transactionSpeedChan <- true
}

func GetDifficultyTarget() int {
	<-difficultyTargetChan
	retVal := difficultyTarget
	difficultyTargetChan <- true
	return retVal
}

func SetMinNeighbours(value int) {
	<-minNeighboursChan
	minNeighbours = value
	minNeighboursChan <- true
}

func SetMaxNeighbours(value int) {
	<-maxNeighboursChan
	maxNeighbours = value
	maxNeighboursChan <- true
}

func SetMaxSecondsPingDelay(value int) {
	<-maxSecondsPingDelayChan
	maxSecondsPingDelay = value
	maxSecondsPingDelayChan <- true
}

func SetMaxTransactionsInBlock(value int) {
	<-maxTransactionsInBlockChan
	maxTransactionsInBlock = value
	maxTransactionsInBlockChan <- true
}

func SetMinTransactionsInBlock(value int) {
	<-minTransactionsInBlockChan
	minTransactionsInBlock = value
	minTransactionsInBlockChan <- true
}

func SetDifficultyTarget(value int) {
	<-difficultyTargetChan
	difficultyTarget = value
	difficultyTargetChan <- true
}
