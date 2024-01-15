package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/bootstrap"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/client"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/mining"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/nodestate"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/propagation"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/server"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

type NodeAddress = types.NodeAddress

func main() {
	var (
		port            int
		bootstrapPort   int
		isBootstrapNode bool
	)

	flag.IntVar(&port, "p", 0, "The server port")
	flag.IntVar(&bootstrapPort, "b", 0, "The bootstrap node's port")
	flag.BoolVar(&isBootstrapNode, "m", false, "If this node is the bootstrap node, set this flag; -b will be used as the bootstrap node's bootstrap port")
	flag.Parse()

	if flag.NFlag() < 2 {
		flag.Usage()
		os.Exit(1)
	}

	// if both -b and -m are set
	if port == 0 || bootstrapPort == 0 {
		// log.Println("Error: Both -p and -b flags must be set, -m flag is optional (if -m is set, it means this is the bootstrap node, and -b will be used as the bootstrap node's bootstrap port)")
		os.Exit(1)
	}

	var serverNode NodeAddress = types.NewNodeAddress(port)
	var bootstrapNode NodeAddress = types.NewNodeAddress(bootstrapPort)

	// Initialize the configuration
	configuration.InitConfiguration(serverNode, bootstrapNode, bootstrapNode, isBootstrapNode)

	// Initialize state of the node
	nodestate.InitNodeState()

	if isBootstrapNode { // if -m is setTransaction
		// This is the bootstrap node of the network
		// log.Println("This is the bootstrap node")
		// log.Println("Bootstrap node's Server port:", port)
		// log.Println("Bootstrap Node's Bootstrap port :", bootstrapPort)

		// Start the bootstrap server
		// log.Println("Starting the bootstrap server")
		go bootstrap.StartBootstrapServer(bootstrapNode)

		// Start the normal server
		// log.Println("Starting the server")
		go server.StartServer()

		// Start the client
		// log.Println("Starting the client")
		go client.StartClient()

	} else { // if -b is set
		// This is a regular node
		// log.Println("This is a regular node")
		// log.Println("Regular node's Server port:", port)
		// log.Println("Regular node's Bootstrap port:", bootstrapPort)

		// Start the normal server
		// log.Println("Starting the server")
		go server.StartServer()

		// Start the client
		// log.Println("Starting the client")
		go client.StartClient()
	}

	// Keep the program running
	// select {}
	// go periodicallyGenerateRandomTransaction()
	Menu()
	os.Exit(0)
}

// func periodicallyDisplayTransactionPool() {
// 	for {
// 		txPool := nodestate.ReadTempTxPool()
// 		txPool.DisplayValueFromPool()
// 		time.Sleep(10000 * time.Millisecond)
// 	}
// }

func displayBlockChain() {
	for {
		bchain := nodestate.ReadBlockChain()
		bchain.Display()
		time.Sleep(1000 * time.Millisecond)
	}
}

func ClearScreen() {
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

func periodicallyGenerateRandomTransaction() {
	counter := 0
	for {
		counter++
		client.CreateRandomTransaction(counter)
		sleepTime := configuration.GetTransactionSpeed()
		// Sleep for 1 seconds
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}
}

func Menu() {
	// Assuming blockchain is initialized and accessible
	// transactions1 := types.TransactionList{
	// 	Transactions: []types.Transaction{
	// 		{Value: "Tx1", Hash: "df1056b304b48a5090d0d43420cea4e8a70bfd9b30757c46647abda7d648e971"},
	// 		{Value: "Tx2", Hash: "df1056b304b48a5090d0d43420cea4e8a70bfd9b30757c46647abda7d648e971"},
	// 		{Value: "Tx3", Hash: "df1056b304b48a5090d0d43420cea4e8a70bfd9b30757c466   47abda7d648e971"},
	// 		{Value: "Tx4", Hash: "df1056b304b48a5090d0d43420cea4e8a70bfd9b30757c46647abda7d648e971"},
	// 		{Value: "Tx5", Hash: "df1056b304b48a5090d0d43420cea4e8a70bfd9b30757c46647abda7d648e971"},
	// 		{Value: "Tx6", Hash: "df1056b304b48a5090d0d43420cea4e8a70bfd9b30757c46647abda7d648e971"},
	// 	},
	// }

	// transactions2 := types.TransactionList{
	// 	Transactions: []types.Transaction{
	// 		{Value: "Tx7", Hash: "df1056b304b48a5090d0d43420cea4e8a70bfd9b30757c46647abda7d648e971"},
	// 		{Value: "Tx8", Hash: "df1056b304b48a5090d0d43420cea4e8a70bfd9b30757c46647abda7d648e971"},
	// 		{Value: "Tx9", Hash: "df1056b304b48a5090d0d43420cea4e8a70bfd9b30757c46647abda7d648e971"},
	// 		{Value: "Tx10", Hash: "df1056b304b48a5090d0d43420cea4e8a70bfd9b30757c46647abda7d648e971"},
	// 		{Value: "Tx11", Hash: "df1056b304b48a5090d0d43420cea4e8a70bfd9b30757c46647abda7d648e971"},
	// 	},
	// }

	// transactions3 := types.TransactionList{
	// 	Transactions: []types.Transaction{
	// 		{Value: "Tx12", Hash: "df1056b304b48a5090d0d43420cea4e8a70bfd9b30757c46647abda7d648e971"},
	// 		{Value: "Tx13", Hash: "df1056b304b48a5090d0d43420cea4e8a70bfd9b30757c46647abda7d648e971"},
	// 		{Value: "Tx14", Hash: "df1056b304b48a5090d0d43420cea4e8a70bfd9b30757c46647abda7d648e971"},
	// 		{Value: "Tx15", Hash: "df1056b304b48a5090d0d43420cea4e8a70bfd9b30757c46647abda7d648e971"},
	// 	},
	// }

	// // Create blocks using the CreateBlock function
	// block1 := types.CreateNewBlock(transactions1, "0000000000000000000000000000000000000000000000000000000000000000", 1, 8)
	// block2 := types.CreateNewBlock(transactions2, block1.BlockHash, 2, 8)
	// block3 := types.CreateNewBlock(transactions3, block2.BlockHash, 3, 8)

	// // Create a blockchain and link the blocks
	// blockchain := types.BlockChain{}
	// node3 := &types.BlockNode{Block: block3, PrevNode: nil}
	// node2 := &types.BlockNode{Block: block2, PrevNode: node3}
	// node1 := &types.BlockNode{Block: block1, PrevNode: node2}
	// blockchain.LatestNode = node1

	// nodestate.SetBlockChain(blockchain)

	var choice int

	for {
		fmt.Printf("Node Address: %s\n", configuration.GetSelfServerAddress().GetAddress())
		fmt.Println("Menu:")
		fmt.Println("1. Display Merkle Tree")
		fmt.Println("2. Display Blockchain")
		fmt.Println("3. ChangeBlock")
		fmt.Println("4. Display Network List") // New option added
		fmt.Println("5. Blockchain Configuration")
		fmt.Println("6. Display Transaction Pool")
		fmt.Println("7. BootStrap Nodes List")
		fmt.Println("8. Create New Transaction")
		fmt.Println("9. Exit")

		fmt.Print("Enter your choice: ")
		fmt.Scanln(&choice)
		ClearScreen()

		switch choice {
		case 1:
			fmt.Println("------------------------------")
			fmt.Println("                       ")
			bchain := nodestate.ReadBlockChain()
			types.DisplayMerkleTree(&bchain)
			fmt.Println("------------------------------")
			fmt.Println("                       ")
		case 2:
			fmt.Println("------------------------------")
			fmt.Println("                       ")
			bchain := nodestate.ReadBlockChain()
			types.DisplayBlockchain(&bchain)
			fmt.Println("------------------------------")
			fmt.Println("                       ")
		case 3:
			fmt.Println("------------------------------")
			fmt.Println("                       ")
			bchain := nodestate.ReadBlockChain()
			changeBlock(&bchain)
			fmt.Println("------------------------------")
			fmt.Println("                       ")
		case 4:
			fmt.Println("------------------------------")
			fmt.Println("                       ")
			nl := propagation.GetP2PNetwork()
			displayNetworkLists(nl) // New function for displaying network lists
			fmt.Println("------------------------------")
			fmt.Println("                       ")
		case 6:
			fmt.Println("------------------------------")
			fmt.Println("                       ")
			txPool := nodestate.ReadTransactionPool()
			txPool.DisplayTransactionPool()
			fmt.Println("------------------------------")
			fmt.Println("                       ")
		case 5:
			fmt.Println("------------------------------")
			fmt.Println("                       ")
			blockchainConfiguration()
			fmt.Println("------------------------------")
			fmt.Println("                       ")
		case 7:
			fmt.Println("------------------------------")
			fmt.Println("                       ")
			_, bootstrapNodesList := nodestate.ReadCurrentExistingNodes("")
			for _, node := range bootstrapNodesList.BootstrapNodes {
				types.PrintBootstrapNodeAddress(node)
			}
			fmt.Println("------------------------------")
			fmt.Println("                       ")
		case 8:
			fmt.Println("------------------------------")
			fmt.Println("                       ")
			client.CreateTransaction()
			fmt.Println("------------------------------")
			fmt.Println("                       ")
		case 9:
			fmt.Println("------------------------------")
			fmt.Println("---------------Exiting---------------")
			return
		default:
			fmt.Println("------------------------------")
			fmt.Println("                       ")
			fmt.Println("Invalid choice. Please select a valid option.")
			fmt.Println("------------------------------")
			fmt.Println("                       ")
		}
	}
}

func displayNetworkLists(nl []types.NetworkList) {
	for _, networkList := range nl {
		fmt.Println("------------------------------")
		// fmt.Println("                       ")
		// fmt.Println("Source Node:")
		// fmt.Println(networkList.SourceNode)
		// fmt.Println("Destination Nodes:")
		// networkList.DestinationNodes.Display()
		// fmt.Println("                       ")
		networkList.Display()
	}

}

func changeBlock(blockchain *types.BlockChain) {
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
	var block *types.Block
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

		newTransaction := types.NewTransaction(newValue)

		// Update the transaction in the TransactionList
		block.Body.Transactions.Transactions[transactionIndex] = newTransaction

	}

	// Recalculate Merkle tree and update Merkle root hash of the block
	newMerkleTree := types.NewMerkleTree(block.Body.Transactions)
	block.Body.MerkleTree = newMerkleTree
	block.Header.MerkleRootHash = newMerkleTree.MerkleRoot()

	// Step 9: Transactions are already added to the block by modifying the slice
	fmt.Println("Transactions updated successfully.")

	block.Display()
	newMerkleTree.Display()

	// Starting mining from the modified block
	mining.HandleTamperedBlockchain(blockchain, blockHeight)
}

func blockchainConfiguration() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nBlockchain Configuration Menu")
		fmt.Println("0. Set Difficulty Target (Current: ", configuration.GetDifficultyTarget())
		fmt.Printf("1. Set Minimum Neighbours (Current: %d)\n", configuration.GetMinNeighbours())
		fmt.Printf("2. Set Maximum Neighbours (Current: %d)\n", configuration.GetMaxNeighbours())
		fmt.Printf("3. Set Max Seconds Ping Delay (Current: %d)\n", configuration.GetMaxSecondsPingDelay())
		fmt.Printf("4. Set Max Transactions In Block (Current: %d)\n", configuration.GetMaxTransactionsInBlock())
		fmt.Printf("5. Set Min Transactions In Block (Current: %d)\n", configuration.GetMinTransactionsInBlock())
		fmt.Printf("7. Set Transaction Speed (Current: %d)\n", configuration.GetTransactionSpeed())
		fmt.Println("6. Return to Main Menu")

		fmt.Print("Enter your choice: ")
		choiceStr, _ := reader.ReadString('\n')
		choice, err := strconv.Atoi(strings.TrimSpace(choiceStr))
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		if choice == 6 {
			break
		}

		fmt.Print("Enter new value: ")
		valueStr, _ := reader.ReadString('\n')
		value, err := strconv.Atoi(strings.TrimSpace(valueStr))
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		switch choice {
		case 0:
			configuration.SetDifficultyTarget(value)
			fmt.Println("Difficulty Target set to:", configuration.GetDifficultyTarget())
		case 1:
			configuration.SetMinNeighbours(value)
			fmt.Println("Minimum Neighbours set to:", configuration.GetMinNeighbours())
		case 2:
			configuration.SetMaxNeighbours(value)
			fmt.Println("Maximum Neighbours set to:", configuration.GetMaxNeighbours())
		case 3:
			configuration.SetMaxSecondsPingDelay(value)
			fmt.Println("Max Seconds Ping Delay set to:", configuration.GetMaxSecondsPingDelay())
		case 4:
			configuration.SetMaxTransactionsInBlock(value)
			fmt.Println("Max Transactions In Block set to:", configuration.GetMaxTransactionsInBlock())
		case 5:
			configuration.SetMinTransactionsInBlock(value)
			fmt.Println("Min Transactions In Block set to:", configuration.GetMinTransactionsInBlock())
		case 7:
			configuration.SetTransactionSpeed(value)
			fmt.Println("Transaction Speed set to:", configuration.GetTransactionSpeed())
		default:
			fmt.Println("Invalid choice. Please select a valid option.")
		}
	}
}
