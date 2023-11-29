package mining

import (
	"os"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/validation"
)

// VerifyBlock checks if the block's hash and Merkle root are valid
func VerifyBlock(b types.Block) types.ReturnType {
	// Verify the mining proof

	// Generate the target hash
	diffTarget := b.Header.DifficultyTarget
	_, targetHashBytes := validation.GenerateTargetHash(diffTarget)
	// Recalculate the block header hash
	_, hashBytes := b.RecalculateBlockHash()
	if !validation.CompareWithTargetHash(hashBytes, targetHashBytes) {
		// fmt.Println("Mining proof verification failed!")
		// fmt.Println("Difficulty target:", diffTarget)
		// fmt.Println("Current Nonce: ", b.Header.Nonce)
		// fmt.Println("Target hash:", targetHashString)
		// fmt.Println("Block hash:", hashString)
		os.Exit(1)
		return types.MiningProofFailed
	}

	// Verify the Merkle root
	if b.IsTransactionListTampered() {
		// fmt.Println("Merkle root verification failed!")
		return types.MerkleRootFailed
	}

	// fmt.Println("Block verification successful!")
	return types.BlockVerificationSuccessful
}

func VerifyBlockChain(bchain types.BlockChain) (types.ReturnType, *types.BlockNode, *types.BlockNode) {
	// Verify the blocks in the blockchain

	prevBlockHash := ""
	prevHeight := -1
	var prevNode *types.BlockNode = nil
	currentNode := bchain.LatestNode

	for currentNode != nil {
		// Verify the current block (mining proof, Merkle root)
		if ret := VerifyBlock(currentNode.Block); ret != types.BlockVerificationSuccessful {
			return ret, currentNode, prevNode
		}

		// Verify the block height
		if currentNode.Block.Header.Height != prevHeight-1 && currentNode != bchain.LatestNode {
			// fmt.Println("Block height verification failed!")
			return types.BlockHeightFailed, currentNode, prevNode
		}

		// Verify the previous block hash
		if currentNode.Block.BlockHash != prevBlockHash && currentNode != bchain.LatestNode {
			// fmt.Println("Previous block hash verification failed!")
			return types.PreviousBlockHashFailed, currentNode, prevNode
		}

		prevHeight = currentNode.Block.Header.Height
		prevBlockHash = currentNode.Block.BlockHash
		prevNode = currentNode
		currentNode = currentNode.PrevNode
	}

	// Make sure that the last node is the genesis block
	if prevNode.Block.Header.Height != 0 {
		// fmt.Println("Genesis block height verification failed!")
		return types.GenesisBlockHeightFailed, currentNode, prevNode
	}

	return types.BlockChainVerificationSuccessful, nil, nil
}

// func VerifyBlockChain(bchain types.BlockChain) (int, bool) {
// 	// Verify the blocks in the blockchain
// 	// Verify the block headers
// 	// Verify the Merkle roots
// 	// Verify the mining proof
// 	// Verify the transactions in the blocks (no duplicate transactions)

// 	lowestHeight, ok := -1, true

// 	// If the blockchain is empty, return true
// 	if bchain.LatestNode == nil {
// 		return lowestHeight, ok
// 	}
// 	currentNode := bchain.LatestNode
// 	prevNode := bchain.LatestNode.PrevNode

// 	for prevNode != nil {
// 		currentHeight := currentNode.Block.Header.Height

// 		// Verify the current block (mining proof, Merkle root)
// 		if !VerifyBlock(currentNode.Block) {
// 			if currentHeight < lowestHeight {
// 				ok = false
// 				lowestHeight = currentHeight
// 			}
// 		}

// 		// Verify the previous block hash
// 		if currentNode.Block.Header.PreviousBlockHash != prevNode.Block.BlockHash {
// 			// fmt.Println("Previous block hash verification failed!")
// 			if currentHeight-1 < lowestHeight {
// 				ok = false
// 				lowestHeight = currentHeight - 1
// 			}
// 		}

// 		// Verify the block height
// 		if currentNode.Block.Header.Height != prevNode.Block.Header.Height+1 {
// 			// fmt.Println("Block height verification failed!")
// 			if currentHeight-1 < lowestHeight {
// 				ok = false
// 				lowestHeight = currentHeight
// 			}
// 		}
// 	}

// 	// Verify the genesis block

// 	// Verify the current block (mining proof, Merkle root)
// 	if !VerifyBlock(currentNode.Block) {
// 		return false
// 	}

// 	// Verify the block height
// 	if currentNode.Block.Header.Height != 0 {
// 		// fmt.Println("Genesis block height verification failed!")
// 		return false
// 	}

// 	// Verify the previous block hash
// 	if currentNode.Block.Header.PreviousBlockHash != "00000000000000000000000000000000" {
// 		// fmt.Println("Genesis block previous hash verification failed!")
// 		return false
// 	}

// 	// Verify the transactions in the blockchain
// 	if _, ok := VerifyDuplicateTransactions(currentNode); !ok {
// 		return false
// 	}

// 	// fmt.Println("Blockchain verification successful!")
// 	return true
// }

func GetMissingBlocks(bchain types.BlockChain) []int {
	// Verify that there are no missing blocks in the blockchain
	// Traverse through the blockchain
	// If the current block height is not equal to the previous block height + 1, return false and the missing block heights
	// Otherwise, return true and an empty slice

	if bchain.LatestNode == nil {
		return nil
	}

	missingBlocks := make([]int, 0)

	prevNode := bchain.LatestNode
	currentNode := bchain.LatestNode.PrevNode

	for currentNode != nil {
		prevHeight := prevNode.Block.Header.Height
		currentHeight := currentNode.Block.Header.Height

		if currentHeight+1 != prevHeight {
			// Check how many blocks are missing
			for i := currentHeight + 1; i < prevHeight; i++ {
				missingBlocks = append(missingBlocks, i)
			}
		}
		prevNode = currentNode
		currentNode = currentNode.PrevNode
	}
	// Check the genesis block
	if prevNode.Block.Header.Height != 0 {
		// Check how many blocks are missing
		for i := 0; i < prevNode.Block.Header.Height; i++ {
			missingBlocks = append(missingBlocks, i)
		}
	}
	if len(missingBlocks) > 0 {
		return missingBlocks
	} else {
		return nil
	}
}

func VerifyDuplicateTransactionsInBlockChain(bchain types.BlockChain) (int, bool) {
	// Verify that there are no duplicate transactions in the block
	// Create a map of transaction hashes
	// If a transaction hash already exists in the map, return false
	// Otherwise, add the transaction hash to the map
	// Return true

	currentNode := bchain.LatestNode
	transactionMap := make(map[string]int)
	for currentNode != nil {
		currentHeight := currentNode.Block.Header.Height
		transactionList := currentNode.Block.Body.Transactions
		for _, tx := range transactionList.Transactions {
			if h, ok := transactionMap[tx.Hash]; ok {
				// fmt.Println("Duplicate transaction found at block height:", h)
				return h, false
			} else {
				transactionMap[tx.Hash] = currentHeight
			}
		}
	}
	return -1, true
}

func VerifyDuplicateTransactionsInBlock(bchain types.BlockChain, b types.Block) bool {
	// Verify that there are no duplicate transactions in the block
	// Create a map of transaction hashes
	// If a transaction hash already exists in the map, return false
	// Otherwise, add the transaction hash to the map
	// Return true

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
				// fmt.Println("Duplicate transaction found in block!")
				return false
			}
		}
	}
	return true
}
