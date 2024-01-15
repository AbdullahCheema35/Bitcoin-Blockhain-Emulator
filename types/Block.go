package types

import (
	"fmt"
	"time"
)

// Block struct containing both header and body
type Block struct {
	BlockHash string
	Header    BlockHeader
	Body      BlockBody
}

// Function to calculate the hash of the entire block
func (b *Block) RecalculateBlockHash() (string, []byte) {
	blockHashString, blockHashBytes := b.Header.RecalculateBlockHeaderHash()
	b.BlockHash = blockHashString
	return blockHashString, blockHashBytes
}

// Function to create a new block
func CreateNewBlock(transactions TransactionList, previousBlockHash string, height int, diffTarget int) Block {
	// Build Merkle tree from transactions
	merkleTree := NewMerkleTree(transactions)

	header := BlockHeader{
		PreviousBlockHash: previousBlockHash,
		Nonce:             0,
		Height:            height,
		MerkleRootHash:    merkleTree.MerkleRoot(),
		Timestamp:         time.Now(),
		DifficultyTarget:  diffTarget,
	}

	body := BlockBody{
		Transactions: transactions,
		MerkleTree:   merkleTree,
	}

	blockHashString, _ := header.RecalculateBlockHeaderHash()

	return Block{
		Header:    header,
		Body:      body,
		BlockHash: blockHashString,
	}
}

// Function to check if the transaction list has been tampered with
func (b Block) IsTransactionListTampered() bool {
	// Recalculate the Merkle root based on the updated transactions
	recalculatedMerkleRoot := NewMerkleTree(b.Body.Transactions).MerkleRoot()

	// Compare the recalculated Merkle root with the one stored in the block header
	return recalculatedMerkleRoot != b.Header.MerkleRootHash
}

func (b Block) DisplayHeader() {
	fmt.Println("Block Header:")
	fmt.Printf("PBH: %s\n", b.Header.PreviousBlockHash)
	// fmt.Printf("Nonce: %d\n", b.Header.Nonce)
	fmt.Printf("Height: %d\n", b.Header.Height)
	fmt.Printf("Block Hash: %s\n", b.BlockHash)
	// fmt.Printf("Merkle Root Hash: %s\n", b.Header.MerkleRootHash)
	// fmt.Printf("Timestamp: %s\n", b.Header.Timestamp.String())
	fmt.Println("------------------------------")
}

// Function to display the block
func (b Block) Display() {
	// fmt.Println("Block Header:")
	// fmt.Printf("Previous Block Hash: %s\n", b.Header.PreviousBlockHash)
	// fmt.Printf("Nonce: %d\n", b.Header.Nonce)
	// fmt.Printf("Height: %d\n", b.Header.Height)
	// // fmt.Printf("Block Hash: %s\n", b.Header.BlockHash)
	// fmt.Printf("Merkle Root Hash: %s\n", b.Header.MerkleRootHash)
	// fmt.Printf("Timestamp: %s\n", b.Header.Timestamp.String())

	fmt.Println("Block Transactions:")
	for _, tx := range b.Body.Transactions.Transactions {
		fmt.Printf("Tx Hash: %s, Tx Value: %s\n", tx.Hash, tx.Value)
	}
	fmt.Println("------------------------------")
	// fmt.Println()
	// b.Body.Transactions.DisplayTransactionPool()
	// fmt.Println("------------------------------")
}

// Function to recalculate Merkle Root Hash based on current block transactions
func (b *Block) RecalculateMerkleRoot() string {
	// Recalculate Merkle Root Hash
	newMerkleTreeRoot := NewMerkleTree(b.Body.Transactions).MerkleRoot()
	b.Header.MerkleRootHash = newMerkleTreeRoot
	return newMerkleTreeRoot
}

// // Function to update block information in case of tampering
// func (b *Block) UpdateBlock(transactions TransactionList) {
// 	// Tamper the first transaction's value
// 	//b.TamperTransaction(2, "Txx")

// 	// Recalculate the Merkle root and update block information
// 	b.Header.Nonce++
// 	b.Header.Timestamp = time.Now()
// 	b.Body.MerkleTree = NewMerkleTree(transactions)
// 	b.Header.MerkleRootHash = b.Body.MerkleTree.MerkleRoot()
// 	b.Body.Transactions = transactions
// 	//b.Header.BlockHash = b.CalculateBlockHash()
// }

// func (b *Block) TamperTransaction(transactionIndex int, newValue string) {
// 	if transactionIndex >= 0 && transactionIndex < len(b.Body.Transactions.Transactions) {
// 		// Save the original hash
// 		//originalHash := b.Body.Transactions.Transactions[transactionIndex].Hash

// 		// Tamper the transaction value
// 		b.Body.Transactions.Transactions[transactionIndex] = newValue

// 		// Calculate new hash for the tampered transaction value
// 		newHash := CalculateTransactionHash(newValue)

// 		// Update the transaction hash
// 		b.Body.Transactions.Transactions[transactionIndex].Hash = newHash

// 		// Recalculate the Merkle root
// 		b.Header.Timestamp = time.Now().Unix()
// 		b.Body.MerkleTree = NewMerkleTree(b.Body.Transactions)
// 		b.Header.MerkleRootHash = b.Body.MerkleTree.MerkleRoot()

// 		// Restore the original hash for consistency
// 		//b.Body.Transactions.Transactions[transactionIndex].Hash = originalHash
// 	}
// }
