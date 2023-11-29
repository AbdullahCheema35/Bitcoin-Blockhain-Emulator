package client

import (
	"time"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/mineblock"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/nodestate"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

func EnoughTransactionsForBlock() bool {
	minTransactions := configuration.GetMinTransactionsInBlock()
	currentTransactions := len(nodestate.ReadTransactionPool().Transactions)
	if currentTransactions >= minTransactions {
		// display the transaction pool
		// fmt.Println("Transaction pool:")
		// for _, transaction := range nodestate.ReadTransactionPool().Transactions {
		// 	fmt.Println(transaction)
		// }
		return true
	}
	return false
}

func purifyTransactionPool(pool types.TransactionList) types.TransactionList {
	// Traverse the blockchain and remove all the transactions that are present in the blockchain
	// from the transaction pool

	bchain := nodestate.ReadBlockChain()
	currentNode := bchain.LatestNode
	for currentNode != nil {
		for _, blockTransaction := range currentNode.Block.Body.Transactions.Transactions {
			pool.RemoveTransaction(blockTransaction)
		}
		currentNode = currentNode.PrevNode
	}
	return pool
}

func createBlock() (bool, types.Block) {
	transactionPool := nodestate.ReadTransactionPool()
	transactionPool = purifyTransactionPool(transactionPool)

	if len(transactionPool.Transactions) < configuration.GetMinTransactionsInBlock() {
		// fmt.Println("Not enough transactions in the Purified Transaction pool to create a block")
		// Print the purified transaction pool
		// for _, transaction := range transactionPool.Transactions {
		// 	fmt.Println(transaction)
		// }
		return false, types.Block{}
	}

	var selectedTransactions []types.Transaction

	if len(transactionPool.Transactions) > configuration.GetMaxTransactionsInBlock() {
		selectedTransactions = transactionPool.GetTransactions()[:configuration.GetMaxTransactionsInBlock()]
	} else {
		selectedTransactions = transactionPool.GetTransactions()
	}

	// Fetch latest block from the blockchain
	bchain := nodestate.LockBlockChain()
	defer func() {
		nodestate.UnlockBlockChain(bchain)
	}()

	latestBlockHash := bchain.GetLatestBlockHash()
	latestBlockHeight := bchain.GetLatestBlockHeight()

	selectedTransactionList := types.NewTransactionListFromSlice(selectedTransactions)

	diffTarget := configuration.GetDifficultyTarget()

	newBlock := types.CreateNewBlock(selectedTransactionList, latestBlockHash, latestBlockHeight+1, diffTarget)
	return true, newBlock
}

func StartCreateBlocks() {
	var miningChan chan bool
	for {
		if EnoughTransactionsForBlock() {
			miningChan = nodestate.InitMiningChan()
			ok, newCreatedBlock := createBlock()
			if ok {
				mineblock.MineBlock(newCreatedBlock, miningChan)
				// fmt.Println("Mine block function returned")
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}
