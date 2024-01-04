package mining

import (
	"log"
	"os"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/nodestate"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

func HandleNewBlock(block types.Block, receivedFrom types.NodeAddress) types.ReturnType {
	result := AddNewBlockToBlockChain(block)

	switch result {
	case types.NewBlockVerificationFailed:
		log.Println("New block verification failed")
		return types.DoNothing
	case types.NewHeightLEQCurrentHeight:
		log.Println("New block height is less than or equal to the current height")
		return types.DoNothing
	case types.NewBlockPrevHashDontMatch:
		log.Println("New block previous hash does not match the hash of the latest block in the blockchain")
		selfNode := configuration.GetSelfServerAddress()
		if selfNode.GetAddress() == receivedFrom.GetAddress() {
			return types.DoNothing
		}
		return types.InitiateBroadcastBlockChainRequest
	case types.NewBlockDuplicateTransactions:
		log.Println("New block contains duplicate transactions")
		return types.DoNothing
	case types.NewBlockAddedSuccessfully:
		log.Println("New block added successfully. Stopping mining and broadcasting the block.")
		AbortTheMiningProcess()
		return types.InitiateBroadcastBlock
	default:
		log.Println("Serious Error: Invalid return type from AddNewBlockToBlockChain")
		return types.DoNothing
	}
}

func HandleNewBChain(newbchain types.BlockChain) {
	newHeight := newbchain.GetLatestBlockHeight()

	bchain := nodestate.LockBlockChain()
	defer func() {
		nodestate.UnlockBlockChain(bchain)
	}()

	currentHeight := bchain.GetLatestBlockHeight()
	if newHeight <= currentHeight {
		log.Println("Received blockchain is not longer than the current blockchain. Do nothing")
		return
	}

	result, _, _ := VerifyBlockChain(newbchain)

	if result == types.BlockChainVerificationSuccessful {
		log.Println("Received blockchain is valid. Replacing the current blockchain and stopping mining.")
		nodestate.SetBlockChain(newbchain)
		AbortTheMiningProcess()
		return
	}
}

func AbortTheMiningProcess() {
	nodestate.CloseMiningChan()
}

func AddNewBlockToBlockChain(b types.Block) types.ReturnType {
	ret := VerifyBlock(b)
	if ret != types.BlockVerificationSuccessful {
		return types.NewBlockVerificationFailed
	}

	newHeight := b.Header.Height

	bchain := nodestate.LockBlockChain()
	defer func() {
		nodestate.UnlockBlockChain(bchain)
	}()

	currentHeight := bchain.GetLatestBlockHeight()
	if newHeight <= currentHeight {
		return types.NewHeightLEQCurrentHeight
	}

	latestBlockHash := bchain.GetLatestBlockHash()
	if b.Header.PreviousBlockHash != latestBlockHash {
		return types.NewBlockPrevHashDontMatch
	}

	verified := VerifyDuplicateTransactionsInBlock(bchain, b)
	if !verified {
		return types.NewBlockDuplicateTransactions
	}

	success := bchain.AddBlock(b)
	if !success {
		log.Println("Serious Error: Block was not added to the blockchain")
		os.Exit(1)
	}
	// log.Println("Transaction Pool Before Pruning:")
	// nodestate.ReadTransactionPool().DisplayTransactionPool()
	PruneTransactionList(b)
	// log.Println("Transaction Pool After Pruning:")
	// nodestate.ReadTransactionPool().DisplayTransactionPool()
	return types.NewBlockAddedSuccessfully
}

func PruneTransactionList(block types.Block) {
	transactionPool := nodestate.LockTransactionPool()
	defer func() {
		nodestate.UnlockTransactionPool(transactionPool)
	}()

	// Keep track of how many transactions were removed
	// prunedTransactions := 0

	for _, transaction := range block.Body.Transactions.Transactions {
		if transactionPool.RemoveTransaction(transaction) {
			// prunedTransactions++
		}
	}
	// log.Printf("Pruned (%d/%d) transactions from the transaction pool\n", prunedTransactions, len(block.Body.Transactions.Transactions))
}
