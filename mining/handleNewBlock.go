package mining

import (
	"os"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/nodestate"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

func HandleNewBlock(block types.Block, receivedFrom types.NodeAddress) types.ReturnType {
	// Handle the received block from a node
	// Verify the received block
	// If the block is valid, add the block to the blockchain, and broadcast the block to all the peers
	// If the block is invalid, discard the block

	result := AddNewBlockToBlockChain(block)

	switch result {
	case types.NewBlockVerificationFailed:
		// log.Println("New block verification failed")
		return types.DoNothing
	case types.NewHeightLEQCurrentHeight:
		// log.Println("New block height is less than or equal to the current height")
		return types.DoNothing
	case types.NewBlockPrevHashDontMatch:
		// log.Println("New block previous hash does not match the hash of the latest block in the blockchain")
		selfNode := configuration.GetSelfServerAddress()
		if selfNode.GetAddress() == receivedFrom.GetAddress() {
			return types.DoNothing
		}
		return types.InitiateBroadcastBlockChainRequest
	case types.NewBlockDuplicateTransactions:
		// selfNode := configuration.GetSelfServerAddress()
		// if selfNode.GetAddress() == receivedFrom.GetAddress() {
		// 	return types.DoNothing
		// }
		// return types.InitiateBroadcastBlockChainRequest
		// log.Println("New block contains duplicate transactions")
		return types.DoNothing
	case types.NewBlockAddedSuccessfully:
		// log.Println("New block added successfully")
		// Stop mining
		AbortTheMiningProcess()
		// Flood the block to all the peers except the one from which the block was received
		return types.InitiateBroadcastBlock
	default:
		// log.Println("\n\n\n\n\n\nSerious Error: Invalid return type from AddNewBlockToBlockChain\n\n\n\n\n\n")
		return types.DoNothing
	}
}

func HandleNewBChain(newbchain types.BlockChain) {
	// Handle the received blockchain from a node
	// Verify the received blockchain
	// If the blockchain is valid, replace the current blockchain with the received blockchain
	// If the blockchain is invalid, discard the blockchain
	newHeight := newbchain.GetLatestBlockHeight()

	// Lock our blockchain
	bchain := nodestate.LockBlockChain()
	defer func() {
		nodestate.UnlockBlockChain(bchain)
	}()

	currentHeight := bchain.GetLatestBlockHeight()
	if newHeight <= currentHeight {
		// log.Println("Received blockchain is not longer than the current blockchain. Do nothing")
		return
	}

	result, _, _ := VerifyBlockChain(newbchain)

	if result == types.BlockChainVerificationSuccessful {
		// Replace the current blockchain with the received blockchain
		nodestate.SetBlockChain(newbchain)
		// Stop mining
		AbortTheMiningProcess()
		// Flood the block to all the peers except the one from which the block was received
		return
	}
}

func AbortTheMiningProcess() {
	// Get the mining process aborted channel
	nodestate.CloseMiningChan()
}

func AddNewBlockToBlockChain(b types.Block) types.ReturnType {
	ret := VerifyBlock(b)
	if ret != types.BlockVerificationSuccessful {
		return types.NewBlockVerificationFailed
	}

	newHeight := b.Header.Height

	// Get and lock the the block chain
	bchain := nodestate.LockBlockChain()
	defer func() {
		nodestate.UnlockBlockChain(bchain)
	}()

	currentHeight := bchain.GetLatestBlockHeight()
	if newHeight <= currentHeight {
		return types.NewHeightLEQCurrentHeight
	}

	// Verify that the previous block hash matches the hash of the latest block in the blockchain
	latestBlockHash := bchain.GetLatestBlockHash()
	if b.Header.PreviousBlockHash != latestBlockHash {
		// Handle this case in the calling function such that the missing blocks are added to the blockchain
		return types.NewBlockPrevHashDontMatch
	}
	// Verify that the transactions in the block are not already present in the blockchain
	verified := VerifyDuplicateTransactionsInBlock(bchain, b)
	if !verified {
		return types.NewBlockDuplicateTransactions
	}
	// Now finally adding new block to the blockchain
	// Add the block to the blockchain
	success := bchain.AddBlock(b)
	if !success {
		// log.Println("\n\n\n\n\n\nSerious Error: Block was not added to the blockchain\n\n\n\n\n\n")
		os.Exit(1)
	}
	// First prune the transactions from the transaction pool
	PruneTransactionList(b)
	// Return success flag
	return types.NewBlockAddedSuccessfully
}

func PruneTransactionList(block types.Block) {
	transactionPool := nodestate.LockTransactionPool()
	defer func() {
		nodestate.UnlockTransactionPool(transactionPool)
	}()

	for _, transaction := range block.Body.Transactions.Transactions {
		transactionPool.RemoveTransaction(transaction)
	}
}
