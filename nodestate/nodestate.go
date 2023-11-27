package nodestate

import "github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"

var (
	currentConnectionsChan   chan types.ConnectionsList   = make(chan types.ConnectionsList, 1)
	currentExistingNodesChan chan types.BootstrapNodesMap = make(chan types.BootstrapNodesMap, 1)
	bootstrapLockChan        chan bool                    = make(chan bool, 1)
	transactionPoolChan      chan types.TransactionList   = make(chan types.TransactionList, 1)
	blockchainChan           chan types.BlockChain        = make(chan types.BlockChain, 1)
	miningChan               chan bool
	topologyRequestChan      chan types.TopologyRequest
)

func InitNodeState() {
	currentConnectionsChan <- types.NewNodeConnectionsList()
	currentExistingNodesChan <- types.NewBootstrapNodesMap()
	bootstrapLockChan <- false
	transactionPoolChan <- types.NewTransactionList()
	blockchainChan <- types.NewBlockChain()
}

func InitTopologyChan() chan types.TopologyRequest {
	topologyRequestChan = make(chan types.TopologyRequest, 10)
	return topologyRequestChan
}

func GetTopologyChan() chan types.TopologyRequest {
	return topologyRequestChan
}

func CloseTopologyChan() {
	close(topologyRequestChan)
}

func InitMiningChan() chan bool {
	CloseMiningChan()
	miningChan = make(chan bool, 1)
	return miningChan
}

func GetMiningChan() chan bool {
	return miningChan
}

func CloseMiningChan() {
	close(miningChan)
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

func LockTransactionPool() types.TransactionList {
	return <-transactionPoolChan
}

func UnlockTransactionPool(transactionPool types.TransactionList) {
	transactionPoolChan <- transactionPool
}

func LockBlockChain() types.BlockChain {
	return <-blockchainChan
}

func UnlockBlockChain(blockchain types.BlockChain) {
	blockchainChan <- blockchain
}

// Reader functions
func ReadTransactionPool() types.TransactionList {
	transactionPool := LockTransactionPool()
	UnlockTransactionPool(transactionPool)
	return transactionPool
}

func ReadBlockChain() types.BlockChain {
	blockchain := LockBlockChain()
	UnlockBlockChain(blockchain)
	return blockchain
}
