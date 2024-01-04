package nodestate

import (
	"log"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

var (
	currentConnectionsChan   chan types.ConnectionsList   = make(chan types.ConnectionsList, 1)
	currentExistingNodesChan chan types.BootstrapNodesMap = make(chan types.BootstrapNodesMap, 1)
	transactionPoolChan      chan types.TransactionList   = make(chan types.TransactionList, 1)
	blockchainChan           chan types.BlockChain        = make(chan types.BlockChain, 1)
	miningChan               chan bool                    = make(chan bool, 1)
	topologyRequestChan      chan types.TopologyRequest   = make(chan types.TopologyRequest, 10)
)

func InitNodeState() {
	currentConnectionsChan <- types.NewNodeConnectionsList()
	currentExistingNodesChan <- types.NewBootstrapNodesMap()
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
	select {
	case <-topologyRequestChan:
		return
	default:
	}
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
	select {
	case <-miningChan:
		return
	default:
	}
	close(miningChan)
}

// func LockBootstrapChan() bool {
// 	return <-bootstrapLockChan
// }

// func UnlockBootstrapChan(connected bool) {
// 	bootstrapLockChan <- connected
// }

// Reader functions for channels
// Non blocking read
func ReadCurrentConnections(line string) (int, types.ConnectionsList) {
	currentNeighbours, currentConnections := LockCurrentConnections("For Reading: configuration.go: 56")
	UnlockCurrentConnections(currentConnections, "configuration.go: 57")
	return currentNeighbours, currentConnections
}

func ReadCurrentExistingNodes(line string) (int, types.BootstrapNodesMap) {
	currentNeighbours, currentExistingNodes := LockCurrentExistingNodes("For reading: configuration.go: 62")
	UnlockCurrentExistingNodes(currentExistingNodes, "configuration.go: 63")
	return currentNeighbours, currentExistingNodes
}

// Reader functions for channels
// Blocking read
func LockCurrentConnections(line string) (int, types.ConnectionsList) {
	// // log.Println("Locking current connections", line)
	currentConnections := <-currentConnectionsChan
	currentNeighbours := len(currentConnections.GetNodeConnections())
	// // // log.Println("Locking current neighbours; value:", currentNeighbours)
	// // // log.Println("Locking current connections; value:", currentConnections.GetNodeConnections())
	return currentNeighbours, currentConnections
}

func LockCurrentExistingNodes(line string) (int, types.BootstrapNodesMap) {
	// // log.Println("Locking current existing nodes", line)
	currentExistingNodes := <-currentExistingNodesChan
	currentNeighbours := currentExistingNodes.GetLength()
	// // // log.Println("Locking current neighbours; value:", currentNeighbours)
	// // // log.Println("Locking current connections; value:", currentConnections.GetNodeConnections())
	return currentNeighbours, currentExistingNodes
}

// Writer functions for channels
// Blocking write (if the channel is full, i.e., above the buffer size)
func UnlockCurrentConnections(currentConnections types.ConnectionsList, line string) {
	// // // log.Println("Unlocking current neighbours; value:", currentNeighbours)
	// // log.Println("Unlocking current connections; value:", currentConnections.GetNodeConnections())
	currentConnectionsChan <- currentConnections
	// // log.Println("Unlocked current resources", line)
}

func UnlockCurrentExistingNodes(currentExistingNodes types.BootstrapNodesMap, line string) {
	// // // log.Println("Unlocking current neighbours; value:", currentNeighbours)
	// // log.Println("Unlocking current existing nodes")
	currentExistingNodesChan <- currentExistingNodes
	// // log.Println("Unlocked current resources", line)
}

func LockTransactionPool() types.TransactionList {
	// // log.Println("Locking transaction pool")
	return <-transactionPoolChan
}

func UnlockTransactionPool(transactionPool types.TransactionList) {
	// // log.Println("Unlocking transaction pool")
	transactionPoolChan <- transactionPool
}

func LockBlockChain() types.BlockChain {
	// // log.Println("Locking blockchain")
	return <-blockchainChan
}

func UnlockBlockChain(blockchain types.BlockChain) {
	// // log.Println("Unlocking blockchain")
	blockchainChan <- blockchain
}

func SetBlockChain(newbchain types.BlockChain) {
	LockBlockChain()
	UnlockBlockChain(newbchain)
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

func AddTransactionToPool(value string) (bool, types.Transaction) {
	transaction := types.NewTransaction(value)

	transactionPool := LockTransactionPool()
	defer func() {
		UnlockTransactionPool(transactionPool)
	}()

	isAdded := transactionPool.AddTransaction(transaction)
	if isAdded {
		log.Printf("Transaction %s added to the pool\n", value)
		// Display current transaction pool
		transactionPool.DisplayTransactionPool()
	} else {
		log.Printf("Transaction %s already exists in the pool\n", value)
	}
	// if isAdded {
	// 	// Display the transaction pool
	// 	// for _, tx := range transactionPool.Transactions {
	// 	// 	fmt.Println(tx)
	// 	// }
	// 	// fmt.Println()
	// }
	return isAdded, transaction
}
