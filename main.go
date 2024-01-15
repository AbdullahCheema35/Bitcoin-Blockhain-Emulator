package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/bootstrap"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/client"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
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
	// menu()
	periodicallyGenerateRandomTransaction()
}

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

func checkTransactionPool() {
	for {
		txPool := nodestate.ReadTransactionPool()
		bchain := nodestate.ReadBlockChain()

		// Traverse the transaction pool and check if any duplicate transactions

		// Create mapping of transaction hash to transaction
		// If a transaction hash already exists in the map, then it is a duplicate transaction
		// If a transaction hash does not exist in the map, then add it to the map

		txMap := make(map[string]types.Transaction)

		for _, tx := range txPool.Transactions {
			if _, ok := txMap[tx.Hash]; ok {
				txPool.DisplayTransactionPool()
				bchain.Display()
				log.Fatalln("Duplicate transaction found in the transaction pool")
			} else {
				txMap[tx.Hash] = tx
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func periodicallyGenerateRandomTransaction() {
	counter := 0
	for {
		counter++
		client.CreateRandomTransaction(counter)
		// Sleep for 1 seconds
		time.Sleep(1000 * time.Millisecond)
	}
}

func Menu() {
	// Assuming blockchain is initialized and accessible

	transactions1 := types.TransactionList{
		Transactions: []types.Transaction{
			{Value: "Tx1", Hash: "df1056b304b48a5090d0d43420cea4e8a70bfd9b30757c46647abda7d648e971"},
			{Value: "Tx2", Hash: "df1056b304b48a5090d0d43420cea4e8a70bfd9b30757c46647abda7d648e971"},
			{Value: "Tx3", Hash: "df1056b304b48a5090d0d43420cea4e8a70bfd9b30757c46647abda7d648e971"},
			{Value: "Tx4", Hash: "df1056b304b48a5090d0d43420cea4e8a70bfd9b30757c46647abda7d648e971"},
			{Value: "Tx5", Hash: "df1056b304b48a5090d0d43420cea4e8a70bfd9b30757c46647abda7d648e971"},
			{Value: "Tx6", Hash: "df1056b304b48a5090d0d43420cea4e8a70bfd9b30757c46647abda7d648e971"},
		},
	}

	transactions2 := types.TransactionList{
		Transactions: []types.Transaction{
			{Value: "Tx7", Hash: "df1056b304b48a5090d0d43420cea4e8a70bfd9b30757c46647abda7d648e971"},
			{Value: "Tx8", Hash: "df1056b304b48a5090d0d43420cea4e8a70bfd9b30757c46647abda7d648e971"},
			{Value: "Tx9", Hash: "df1056b304b48a5090d0d43420cea4e8a70bfd9b30757c46647abda7d648e971"},
			{Value: "Tx10", Hash: "df1056b304b48a5090d0d43420cea4e8a70bfd9b30757c46647abda7d648e971"},
			{Value: "Tx11", Hash: "df1056b304b48a5090d0d43420cea4e8a70bfd9b30757c46647abda7d648e971"},
		},
	}

	transactions3 := types.TransactionList{
		Transactions: []types.Transaction{
			{Value: "Tx12", Hash: "df1056b304b48a5090d0d43420cea4e8a70bfd9b30757c46647abda7d648e971"},
			{Value: "Tx13", Hash: "df1056b304b48a5090d0d43420cea4e8a70bfd9b30757c46647abda7d648e971"},
			{Value: "Tx14", Hash: "df1056b304b48a5090d0d43420cea4e8a70bfd9b30757c46647abda7d648e971"},
			{Value: "Tx15", Hash: "df1056b304b48a5090d0d43420cea4e8a70bfd9b30757c46647abda7d648e971"},
		},
	}

	// Create blocks using the CreateBlock function
	block1 := types.CreateNewBlock(transactions1, "0000000000000000000000000000000000000000000000000000000000000000", 1, 8)
	block2 := types.CreateNewBlock(transactions2, block1.BlockHash, 2, 8)
	block3 := types.CreateNewBlock(transactions3, block2.BlockHash, 3, 8)

	// Create a blockchain and link the blocks
	blockchain := types.BlockChain{}
	node3 := &types.BlockNode{Block: block3, PrevNode: nil}
	node2 := &types.BlockNode{Block: block2, PrevNode: node3}
	node1 := &types.BlockNode{Block: block1, PrevNode: node2}
	blockchain.LatestNode = node1

	nodestate.SetBlockChain(blockchain)

	var choice int

	counter := 0

	for {
		fmt.Println("Menu:")
		fmt.Println("1. Display Merkle Tree")
		fmt.Println("2. Display Blockchain")
		fmt.Println("3. ChangeBlock")
		fmt.Println("4. Display Network List") // New option added
		fmt.Println("6. Display Transaction Pool")
		fmt.Println("5. Blockchain Configuration")
		fmt.Println("7. BootStrap Nodes List")
		fmt.Println("8. Exit")

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
			types.ChangeBlock(&bchain)
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
			// blockchainConfiguration()
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
		default:
			fmt.Println("------------------------------")
			fmt.Println("                       ")
			fmt.Println("Invalid choice. Please select a valid option.")
			fmt.Println("------------------------------")
			fmt.Println("                       ")
		}
		counter++
		client.CreateRandomTransaction(counter)
	}
}

func displayNetworkLists(nl []types.NetworkList) {
	for _, networkList := range nl {
		// fmt.Println("------------------------------")
		// fmt.Println("                       ")
		// fmt.Println("Source Node:")
		// fmt.Println(networkList.SourceNode)
		// fmt.Println("Destination Nodes:")
		// networkList.DestinationNodes.Display()
		// fmt.Println("                       ")
		networkList.Display()
	}
}
