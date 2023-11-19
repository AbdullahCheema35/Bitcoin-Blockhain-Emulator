package server

import (
	"encoding/gob"
	"log"
	"net"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

type Transaction = types.Transaction

type NodeAddress = types.NodeAddress

func handleConnection(conn net.Conn, messages chan Transaction) {
	defer conn.Close()

	decoder := gob.NewDecoder(conn)
	var msg Transaction

	err := decoder.Decode(&msg)
	if err != nil {
		log.Println("Error decoding message:", err)
		return
	}

	messages <- msg
}

func StartServer(serverNode NodeAddress) {
	serverAddress := serverNode.GetAddress()
	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	log.Println("Server Node listening on port", serverNode.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn, messages)
	}
}
