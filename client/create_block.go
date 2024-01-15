package client

import (
	"log"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/nodestate"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

func EnoughTransactionsForBlock() bool {
	minTransactions := configuration.GetMinTransactionsInBlock()
	currentTransactions := len(nodestate.ReadTransactionPool().Transactions)
	if currentTransactions >= minTransactions {
		log.Printf("Sufficient transactions (%d) for a new block\n", currentTransactions)
		return true
	}
	return false
}

func purifyTransactionPool(pool types.TransactionList) types.TransactionList {
	// Remove transactions present in the blockchain from the transaction pool

	// Keep track of how many transactions were removed
	purifiedTransactions := 0

	bchain := nodestate.ReadBlockChain()
	currentNode := bchain.LatestNode
	for currentNode != nil {
		for _, blockTransaction := range currentNode.Block.Body.Transactions.Transactions {
			if pool.RemoveTransaction(blockTransaction) {
				purifiedTransactions++
			}
		}
		currentNode = currentNode.PrevNode
	}
	log.Printf("Purified (%d) transactions from the transaction pool\n", purifiedTransactions)
	return pool
}

func createBlock() (bool, types.Block) {
	transactionPool := nodestate.ReadTransactionPool()
	transactionPool = purifyTransactionPool(transactionPool)

	if len(transactionPool.Transactions) < configuration.GetMinTransactionsInBlock() {
		log.Printf("Not enough transactions (%d) in the purified pool to create a block\n", len(transactionPool.Transactions))
		return false, types.Block{}
	}

	var selectedTransactions []types.Transaction

	if len(transactionPool.Transactions) > configuration.GetMaxTransactionsInBlock() {
		selectedTransactions = transactionPool.GetTransactions()[:configuration.GetMaxTransactionsInBlock()]
	} else {
		selectedTransactions = transactionPool.GetTransactions()
	}

	// //Print Seleceted Transactions
	// fmt.Println("Selected Transactions:")
	// transactions := transactionPool.GetTransactions()

	// if len(transactions) == 0 {
	// 	fmt.Println("No transactions in the pool.")
	// }

	// for i, transaction := range transactions {
	// 	fmt.Printf("%d. Hash: %s, Value: %s\n", i+1, transaction.Hash, transaction.Value)
	// }
	// fmt.Println("------------------------")

	bchain := nodestate.LockBlockChain()
	defer func() {
		nodestate.UnlockBlockChain(bchain)
	}()

	latestBlockHash := bchain.GetLatestBlockHash()
	latestBlockHeight := bchain.GetLatestBlockHeight()

	selectedTransactionList := types.NewTransactionListFromSlice(selectedTransactions)

	diffTarget := configuration.GetDifficultyTarget()

	newBlock := types.CreateNewBlock(selectedTransactionList, latestBlockHash, latestBlockHeight+1, diffTarget)
	log.Println("New block created successfully")
	return true, newBlock
}
