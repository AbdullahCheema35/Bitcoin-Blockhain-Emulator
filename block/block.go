package block

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// Transaction struct representing a single transaction
type Transaction struct {
	Value string
	Hash  string // Assuming Hash is a string for simplicity
}

// TransactionList struct representing a list of transactions
type TransactionList struct {
	Transactions []Transaction
}

// BlockHeader struct representing the header of a block
type BlockHeader struct {
	PreviousBlockHash string
	Nonce             int
	Height            int
	BlockHash         string
	MerkleRootHash    string
	Timestamp         int64
}

// BlockBody struct representing the body of a block
type BlockBody struct {
	Transactions TransactionList
	MerkleTree   *MerkleTree // Reference to the Merkle tree
}

// Block struct containing both header and body
type Block struct {
	Header BlockHeader
	Body   BlockBody
}

// Function to calculate the hash of the entire block
func (b *Block) CalculateBlockHash() string {
	headerHash := b.Header.CalculateHeaderHash()
	merkleRootHash := b.Header.MerkleRootHash

	combinedHash := headerHash + merkleRootHash

	hash := sha256.Sum256([]byte(combinedHash))
	return hex.EncodeToString(hash[:])
}

// Function to create a new block
func CreateBlock(transactions TransactionList, previousBlockHash string, height int) Block {
	timestamp := time.Now().Unix()

	// Build Merkle tree from transactions
	merkleTree := NewMerkleTree(transactions)

	header := BlockHeader{
		PreviousBlockHash: previousBlockHash,
		Nonce:             0,
		Height:            height,
		MerkleRootHash:    merkleTree.MerkleRoot(),
		Timestamp:         timestamp,
	}

	// Calculate header hash
	headerHash := header.CalculateHeaderHash()

	// Set the calculated header hash as the block hash
	header.BlockHash = headerHash

	body := BlockBody{
		Transactions: transactions,
		MerkleTree:   merkleTree,
	}

	return Block{
		Header: header,
		Body:   body,
	}
}

// Function to check if the transaction list has been tampered with
func (b *Block) IsTransactionListTampered() bool {
	// Recalculate the Merkle root based on the updated transactions
	recalculatedMerkleRoot := b.Body.MerkleTree.MerkleRoot()

	// Compare the recalculated Merkle root with the one stored in the block header
	return recalculatedMerkleRoot != b.Header.MerkleRootHash
}

// Function to update block information in case of tampering
func (b *Block) UpdateBlock(transactions TransactionList) {
	// Tamper the first transaction's value
	b.TamperTransaction(2, "Txx")

	// Recalculate the Merkle root and update block information
	b.Header.Nonce++
	b.Header.Timestamp = time.Now().Unix()
	b.Body.MerkleTree = NewMerkleTree(transactions)
	b.Header.MerkleRootHash = b.Body.MerkleTree.MerkleRoot()
	b.Body.Transactions = transactions
	b.Header.BlockHash = b.CalculateBlockHash()
}

// Function to display the block
// Function to display the block
func (b *Block) Display() {
	fmt.Println("Block Header:")
	fmt.Printf("Previous Block Hash: %s\n", b.Header.PreviousBlockHash)
	fmt.Printf("Nonce: %d\n", b.Header.Nonce)
	fmt.Printf("Height: %d\n", b.Header.Height)
	fmt.Printf("Block Hash: %s\n", b.Header.BlockHash)
	fmt.Printf("Merkle Root Hash: %s\n", b.Header.MerkleRootHash)
	fmt.Printf("Timestamp: %d\n", b.Header.Timestamp)

	fmt.Println("\nBlock Transactions:")
	for _, tx := range b.Body.Transactions.Transactions {
		fmt.Printf("Transaction Value: %s\n", tx.Value)
		fmt.Printf("Transaction Hash: %s\n", tx.Hash) // Update to display the new hash
	}
	fmt.Println()
}

// Function to calculate the hash of block header
func (bh *BlockHeader) CalculateHeaderHash() string {
	hashInput := bh.PreviousBlockHash + fmt.Sprintf("%d", bh.Nonce) + fmt.Sprintf("%d", bh.Timestamp)
	hash := sha256.Sum256([]byte(hashInput))
	return hex.EncodeToString(hash[:])

}

// Function to calculate the hash of a transaction value
func CalculateHash(value string) string {
	hash := sha256.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}

func (b *Block) TamperTransaction(transactionIndex int, newValue string) {
	if transactionIndex >= 0 && transactionIndex < len(b.Body.Transactions.Transactions) {
		// Save the original hash
		//originalHash := b.Body.Transactions.Transactions[transactionIndex].Hash

		// Tamper the transaction value
		b.Body.Transactions.Transactions[transactionIndex].Value = newValue

		// Calculate new hash for the tampered transaction value
		newHash := CalculateHash(newValue)

		// Update the transaction hash
		b.Body.Transactions.Transactions[transactionIndex].Hash = newHash

		// Recalculate the Merkle root
		b.Header.Timestamp = time.Now().Unix()
		b.Body.MerkleTree = NewMerkleTree(b.Body.Transactions)
		b.Header.MerkleRootHash = b.Body.MerkleTree.MerkleRoot()

		// Restore the original hash for consistency
		//b.Body.Transactions.Transactions[transactionIndex].Hash = originalHash
	}
}

// Function to recalculate Merkle Root Hash based on current block transactions
func (b *Block) RecalculateMerkleRoot() string {
	// Recalculate Merkle Root Hash
	b.Body.MerkleTree = NewMerkleTree(b.Body.Transactions)
	return b.Body.MerkleTree.MerkleRoot()
}
