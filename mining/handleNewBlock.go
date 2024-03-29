package mining

import (
	"fmt"
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
		// log.Println("New block verification failed; do nothing")
		return types.DoNothing
	case types.NewHeightLEQCurrentHeight:
		// log.Println("New block height is less than or equal to the current height; do nothing")
		return types.DoNothing
	case types.NewBlockPrevHashDontMatch:
		// log.Println("New block previous hash does not match the hash of the latest block in the blockchain; InitiateBroadcastBlockChainRequest")
		selfNode := configuration.GetSelfServerAddress()
		if selfNode.GetAddress() == receivedFrom.GetAddress() {
			// log.Panicln("Serious Error: Received invalid block from self node")
			return types.DoNothing
		}
		return types.InitiateBroadcastBlockChainRequest
	case types.NewBlockDuplicateTransactions:
		// log.Println("New block contains duplicate transactions; do nothing")
		return types.DoNothing
	case types.NewBlockAddedSuccessfully:
		// log.Println("New block added successfully. Stopping mining and broadcasting the block; InitiateBroadcastBlock")
		AbortTheMiningProcess()
		return types.InitiateBroadcastBlock
	default:
		// log.Println("Serious Error: Invalid return type from AddNewBlockToBlockChain")
		return types.DoNothing
	}
}

func HandleNewBChain(newbchain types.BlockChain) {
	newHeight := newbchain.GetLatestBlockHeight()

	bchain := nodestate.LockBlockChain()

	// log.Println("--------------------------------------------Displaying Current Blockchain---------------------------------------------")

	// bchain.DisplayHeaderInfo()

	// log.Println("--------------------------------------------Displaying Received Blockchain---------------------------------------------")

	// newbchain.DisplayHeaderInfo()

	// log.Println("---------------------------------------------------------------------------------------------------------------------")

	currentHeight := bchain.GetLatestBlockHeight()
	if newHeight <= currentHeight {
		// log.Println("Received blockchain is not longer than the current blockchain. Do nothing")
		return
	}

	result, _, _ := VerifyBlockChain(newbchain)

	// Unlock the blockchain now
	nodestate.UnlockBlockChain(bchain)

	if result == types.BlockChainVerificationSuccessful {
		// log.Println("Received blockchain is valid. Replacing the current blockchain and stopping mining.")
		nodestate.SetBlockChain(newbchain)
		AbortTheMiningProcess()
	}
}

func HandleTamperedBlockchain(newbchain *types.BlockChain, tamperedBlockHeight int) {
	result, _, _ := VerifyBlockChain(*newbchain)

	if result != types.BlockChainVerificationSuccessful {
		fmt.Println("Blockchain Verification Failed. Restarting mining from tampered block height")

		// Get current Bchain
		currBchain := nodestate.LockBlockChain()

		// Get to node at the tampered block height
		currNode := currBchain.LatestNode

		for currNode != nil && currNode.Block.Header.Height >= tamperedBlockHeight {
			currNode = currNode.PrevNode
		}
		currBchain.LatestNode = currNode

		nodestate.UnlockBlockChain(currBchain)

	} else {
		fmt.Println("BlockChain Verification Successful")
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
