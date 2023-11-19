package main

import (
	"flag"
	"log"
	"os"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/bootstrap"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/client"
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
		log.Println("Error: Both -p and -b flags must be set, -m flag is optional (if -m is set, it means this is the bootstrap node, and -b will be used as the bootstrap node's bootstrap port)")
		os.Exit(1)
	}

	if isBootstrapNode { // if -m is setTransaction
		// This is the bootstrap node of the network
		log.Println("This is the bootstrap node")
		log.Println("Bootstrap node's Server port:", port)
		log.Println("Bootstrap Node's Bootstrap port :", bootstrapPort)

		// Start the normal server
		var serverNode NodeAddress = types.NewNodeAddress(port)
		go server.StartServer(serverNode)

		// Start the bootstrap server
		var bootstrapNode NodeAddress = types.NewNodeAddress(bootstrapPort)
		go bootstrap.StartBootstrapServer(bootstrapNode, serverNode)

		// Start the client
		go client.StartClient(serverNode, bootstrapNode)

	} else { // if -b is set
		// This is a regular node
		log.Println("This is a regular node")
		log.Println("Regular node's Server port:", port)
		log.Println("Regular node's Bootstrap port:", bootstrapPort)

		// Start the normal server
		var serverNode NodeAddress = types.NewNodeAddress(port)
		go server.StartServer(serverNode)

		// Start the client
		var bootstrapNode NodeAddress = types.NewNodeAddress(bootstrapPort)
		go client.StartClient(serverNode, bootstrapNode)
	}

	// Keep the program running
	select {}
}
