// main.go
package main

import (
	"Bitcoin-Blockhain-Emulator/block"
	"fmt"
)

// Function to generate transactions
func main() {
	// Create a block with transactions for the first block
	transactionsBlock1 := block.TransactionList{
		Transactions: []block.Transaction{
			{Value: "Tx1", Hash: block.CalculateHash("Tx1")},
			{Value: "Tx2", Hash: block.CalculateHash("Tx2")},
			{Value: "Tx3", Hash: block.CalculateHash("Tx3")},
			{Value: "Tx4", Hash: block.CalculateHash("Tx4")},
			{Value: "Tx5", Hash: block.CalculateHash("Tx5")},
		},
	}

	// Create the first block
	previousBlockHash := "000000000000" // Replace with an actual hash
	blockHeight := 1
	firstBlock := block.CreateBlock(transactionsBlock1, previousBlockHash, blockHeight)

	// Display original block details for the first block
	fmt.Println("First Block Details:")
	firstBlock.Display()

	// Create transactions for the second block
	transactionsBlock2 := block.TransactionList{
		Transactions: []block.Transaction{
			{Value: "Tx6", Hash: block.CalculateHash("Tx6")},
			{Value: "Tx7", Hash: block.CalculateHash("Tx7")},
			{Value: "Tx8", Hash: block.CalculateHash("Tx8")},
			{Value: "Tx9", Hash: block.CalculateHash("Tx9")},
			//{Value: "Tx10", Hash: block.CalculateHash("Tx10")},
		},
	}

	// Create the second block with the previous hash pointing to the first block's hash
	secondBlock := block.CreateBlock(transactionsBlock2, firstBlock.Header.BlockHash, blockHeight+1)
	secondBlock.Display()
	secondBlock.Body.MerkleTree.Display()

	// Tamper the third transaction value in the first block

	firstBlock.UpdateBlock(transactionsBlock1)

	// Display block details after tampering the first block
	fmt.Println("\nFirst Block After Tampering:")
	firstBlock.Display()
	firstBlock.Body.MerkleTree.Display()

	// Check if the second block's previous hash matches the hash of the first block
	if secondBlock.Header.PreviousBlockHash != firstBlock.Header.BlockHash {
		fmt.Println("\nBlock is tempered!")
	} else {
		fmt.Println("\nBlock is valid.")
	}
}
