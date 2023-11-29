package connection

import (
	"net"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/comm"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/propagation"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

type (
	ConnectionsList = types.ConnectionsList
	NodeConnection  = types.NodeConnection
)

// Establishes a connection to the node at the given address and returns the connection or nil if the connection could not be established
func connectToNode(node types.NodeAddress) net.Conn {
	nodeAddress := node.GetAddress()
	conn, err := net.Dial("tcp", nodeAddress)
	if err != nil {
		// log.Println("Error connecting to node:", err)
		return nil
	}
	return conn
}

func sendConnectionRequestToNode(node types.NodeAddress, conn net.Conn) bool {
	messageType := types.MessageTypeConnectionRequest
	sender := configuration.GetSelfServerAddress()
	messageHeader := types.NewMessageHeader(messageType, sender)
	message := types.NewMessage(messageHeader, nil)

	isMessageSent := comm.SendMessage(conn, message)
	return isMessageSent
}

func receiveConnectionResponseFromNode(conn net.Conn) bool {
	isMessageReceived, message := comm.ReceiveMessage(conn)
	if !isMessageReceived {
		return false
	}
	switch message.Header.Type {
	case types.MessageTypeConnectionResponse:
		// // log.Println("Received a successful connection response from", message.Header.Sender.GetAddress())
		return true
	default:
		// // log.Println("Received an unknown message from", message.Header.Sender.GetAddress())
		conn.Close()
		return false
	}
}

func listenForMessages(nc types.NodeConnection) {
	conn := nc.Conn

	// log.Println("Listening for messages from", nc.Node.GetAddress())

	for {
		err, message := comm.ReceiveMessage(conn)
		if !err {
			// It means connection is broken/lost
			conn.Close()
			HandleLostNodeConnection(nc)
			break
		}
		switch message.Header.Type {
		case types.MessageTypeTransaction:
			sender := message.Header.Sender
			body := message.Body.(types.Transaction)
			propagation.HandleReceivedTransaction(body, sender)
			// log.Printf("Received transaction %v from %v\n", body, sender.GetAddress())
		case types.MessageTypeBlock:
			sender := message.Header.Sender
			body := message.Body.(types.Block)
			propagation.HandleReceivedBlock(body, sender)
			// log.Println("Received a block from", sender.GetAddress())
		case types.MessageTypeBlockChainRequest:
			sender := message.Header.Sender
			// log.Println("Received a blockchain request from", sender.GetAddress())
			propagation.HandleBlockChainRequest(sender)
		case types.MessageTypeBlockChainResponse:
			sender := message.Header.Sender
			// log.Println("Received a blockchain response from", sender.GetAddress())
			propagation.HandleReceivedBlockChain(message.Body.(types.BlockChain), sender)
		case types.MessageTypeTopologyRequest:
			sender := message.Header.Sender
			// log.Println("Received a topology request from", sender.GetAddress())
			propagation.HandleTopologyRequest(sender, message.Body.(types.TopologyRequest))
		case types.MessageTypeTopologyResponse:
			sender := message.Header.Sender
			// log.Println("Received a topology response from", sender.GetAddress())
			propagation.HandleTopologyResponse(message.Body.(types.TopologyRequest), sender)
		default:
			// sender := message.Header.Sender
			// log.Println("Received an unknown message from", sender.GetAddress())
		}
	}
}
