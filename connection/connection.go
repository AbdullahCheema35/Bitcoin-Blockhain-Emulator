package connection

import (
	"log"
	"net"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/common"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
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
		log.Println("Error connecting to node:", err)
		return nil
	}
	return conn
}

func sendConnectionRequestToNode(node types.NodeAddress, conn net.Conn) bool {
	messageType := types.MessageTypeConnectionRequest
	sender := configuration.GetSelfServerAddress()
	messageHeader := types.NewMessageHeader(messageType, sender)
	messageBody := types.ConnectionRequestTypeSuccess
	message := types.NewMessage(messageHeader, messageBody)

	isMessageSent := common.SendMessage(conn, message)
	return isMessageSent
}

func receiveConnectionResponseFromNode(conn net.Conn) bool {
	isMessageReceived, message := common.ReceiveMessage(conn)
	if !isMessageReceived {
		return false
	}
	switch message.Header.Type {
	case types.MessageTypeConnectionResponse:
		if message.Body != types.ConnectionResponseTypeSuccess {
			// log.Printf("Node %v sent a connection response with non-success type\n", message.Header.Sender.GetAddress())
			conn.Close()
			return false
		}
		// log.Println("Received a successful connection response from", message.Header.Sender.GetAddress())
		return true
	default:
		// log.Println("Received an unknown message from", message.Header.Sender.GetAddress())
		conn.Close()
		return false
	}
}

func listenForMessages(nc types.NodeConnection) {
	conn := nc.Conn

	log.Println("Listening for messages from", nc.Node.GetAddress())

	for {
		err, message := common.ReceiveMessage(conn)
		if !err {
			// It means connection is broken/lost
			conn.Close()
			HandleLostNodeConnection(nc)
			break
		}
		switch message.Header.Type {
		case types.MessageTypeTransaction:
			sender := message.Header.Sender
			// TODO: Handle the transaction
			// Temp fix
			transactionData := message.Body.(string)
			log.Printf("Received transaction %v from %v\n", transactionData, sender.GetAddress())
			// End of temp fix
		case types.MessageTypeBlock:
			sender := message.Header.Sender
			log.Println("Received a block from", sender.GetAddress())
		default:
			sender := message.Header.Sender
			log.Println("Received an unknown message from", sender.GetAddress())
		}
	}
}
