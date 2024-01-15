package types

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// type NodeAddress = types.NodeAddress
func ClearScreen() {
	fmt.Println(".......Waiting.........")
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls") // For Windows
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "linux", "darwin":
		cmd := exec.Command("clear") // For Linux/Unix/MacOS
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		fmt.Println("Unsupported platform")
	}
}

func DisplayBlockchain(chain *BlockChain) {
	// ClearScreen()
	currentNode := chain.LatestNode

	for currentNode != nil {
		currentNode.Block.Display()
		currentNode = currentNode.PrevNode
	}
}

func DisplayMerkleTree(chain *BlockChain) {
	ClearScreen()
	DisplayBlockchain(chain) // Display the entire blockchain

	// Ask the user to choose a block by height
	for {
		// Ask the user to choose a block by height
		fmt.Println("Choose a block from the blockchain by entering its height:")
		var chosenHeight int
		_, err := fmt.Scanln(&chosenHeight)
		if err != nil {
			fmt.Println("Invalid input. Please enter a valid integer height.")
			// Clear input buffer
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				break
			}
			continue
		}

		// Find the block by height
		currentNode := chain.LatestNode
		for currentNode != nil {
			if currentNode.Block.Header.Height == chosenHeight {
				fmt.Println("\n\nBlock at chosen height found!")
				currentNode.Block.Display() // Display the chosen block
				fmt.Println()
				currentNode.Block.Body.MerkleTree.Display() // Display the Merkle Tree of the chosen block
				fmt.Println()
				return
			}
			currentNode = currentNode.PrevNode
		}
	}
}

func ChangeBlock(blockchain *BlockChain) {
	reader := bufio.NewReader(os.Stdin)

	// Step 1: Ask the user to enter the block height
	fmt.Print("Enter block height: ")
	blockHeightStr, _ := reader.ReadString('\n')
	blockHeight, err := strconv.Atoi(strings.TrimSpace(blockHeightStr))
	if err != nil {
		fmt.Println("Invalid block height.")
		return
	}

	// Step 2: Iterate through the blockchain to find the block with this height
	var block *Block
	currentNode := blockchain.LatestNode
	for currentNode != nil {
		if currentNode.Block.Header.Height == blockHeight {
			block = &currentNode.Block
			break
		}
		currentNode = currentNode.PrevNode
	}

	if block == nil {
		fmt.Println("Block not found.")
		return
	}

	// Step 3: Display the block
	block.Display()
	block.Body.MerkleTree.Display()

	// Step 4: Ask how many transactions to change
	fmt.Print("How many transactions do you want to change? ")
	transactionCountStr, _ := reader.ReadString('\n')
	transactionCount, err := strconv.Atoi(strings.TrimSpace(transactionCountStr))
	if err != nil {
		fmt.Println("Invalid transaction count.")
		return
	}

	// Step 5: Check if enough transactions exist
	if len(block.Body.Transactions.Transactions) < transactionCount {
		fmt.Println("Not enough transactions in the block.")
		return
	}

	// Steps 6 to 8: Process each transaction change
	for i := 0; i < transactionCount; i++ {
		fmt.Printf("Enter the index of transaction #%d to change: ", i+1)
		transactionIndexStr, _ := reader.ReadString('\n')
		transactionIndex, err := strconv.Atoi(strings.TrimSpace(transactionIndexStr))
		if err != nil || transactionIndex < 0 || transactionIndex >= len(block.Body.Transactions.Transactions) {
			fmt.Println("Invalid transaction index.")
			i-- // Retry the same transaction
			continue
		}

		fmt.Printf("Enter new value for transaction #%d: ", i+1)
		newValue, _ := reader.ReadString('\n')
		newValue = strings.TrimSpace(newValue) // Remove leading/trailing whitespaces

		newTransaction := NewTransaction(newValue)

		// Update the transaction in the TransactionList
		block.Body.Transactions.Transactions[transactionIndex] = newTransaction

	}

	// Recalculate Merkle tree and update Merkle root hash of the block
	newMerkleTree := NewMerkleTree(block.Body.Transactions)
	block.Body.MerkleTree = newMerkleTree
	block.Header.MerkleRootHash = newMerkleTree.MerkleRoot()

	// Step 9: Transactions are already added to the block by modifying the slice
	fmt.Println("Transactions updated successfully.")

	block.Display()
	newMerkleTree.Display()
}
