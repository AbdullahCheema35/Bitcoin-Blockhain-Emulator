package mining

import (
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/validation"
)

func VerifyBlock(b types.Block) types.ReturnType {
	diffTarget := b.Header.DifficultyTarget
	_, targetHashBytes := validation.GenerateTargetHash(diffTarget)

	_, hashBytes := b.RecalculateBlockHash()
	if !validation.CompareWithTargetHash(hashBytes, targetHashBytes) {
		// log.Println("Mining proof verification failed!")
		return types.MiningProofFailed
	}

	if b.IsTransactionListTampered() {
		// log.Println("Merkle root verification failed!")
		return types.MerkleRootFailed
	}

	// log.Println("Block verification successful!")
	return types.BlockVerificationSuccessful
}

func VerifyBlockChain(bchain types.BlockChain) (types.ReturnType, *types.BlockNode, *types.BlockNode) {
	prevBlockHash := ""
	prevHeight := -1
	var prevNode *types.BlockNode = nil
	currentNode := bchain.LatestNode

	for currentNode != nil {
		if ret := VerifyBlock(currentNode.Block); ret != types.BlockVerificationSuccessful {
			return ret, currentNode, prevNode
		}

		if currentNode.Block.Header.Height != prevHeight-1 && currentNode != bchain.LatestNode {
			// log.Println("Block height verification failed!")
			return types.BlockHeightFailed, currentNode, prevNode
		}

		if currentNode.Block.BlockHash != prevBlockHash && currentNode != bchain.LatestNode {
			// log.Println("Previous block hash verification failed!")
			return types.PreviousBlockHashFailed, currentNode, prevNode
		}

		prevHeight = currentNode.Block.Header.Height
		prevBlockHash = currentNode.Block.Header.PreviousBlockHash
		prevNode = currentNode
		currentNode = currentNode.PrevNode
	}

	if prevNode.Block.Header.Height != 0 {
		// log.Println("Genesis block height verification failed!")
		return types.GenesisBlockHeightFailed, currentNode, prevNode
	}

	return types.BlockChainVerificationSuccessful, nil, nil
}

func VerifyDuplicateTransactionsInBlock(bchain types.BlockChain, b types.Block) bool {
	transactionMap := make(map[string]bool)
	transactionList := b.Body.Transactions

	for _, tx := range transactionList.Transactions {
		transactionMap[tx.Hash] = true
	}

	currentNode := bchain.LatestNode

	for currentNode != nil {
		transactionList := currentNode.Block.Body.Transactions

		for _, tx := range transactionList.Transactions {
			if _, ok := transactionMap[tx.Hash]; ok {
				// fmt.Printf("Duplicate transaction (%s) found in the blockchain\n", tx.Value)
				return false
			}
		}
		currentNode = currentNode.PrevNode
	}
	return true
}
