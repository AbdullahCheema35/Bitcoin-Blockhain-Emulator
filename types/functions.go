package types

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
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
