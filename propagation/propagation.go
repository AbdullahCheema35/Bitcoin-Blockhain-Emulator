package propagation

import (
	"log"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/comm"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

// func createTransactionMessage() types.Message {
// 	// Transaction body is randomly constructed string for now
// 	// TODO: Make a transaction type
// 	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
// 	length := 10
// 	b := make([]rune, length)
// 	for i := range b {
// 		b[i] = letters[rand.Intn(len(letters))]
// 	}
// 	transactionData := string(b)
// 	// End of random string generation
// 	// Temp fix

// 	messageType := types.MessageTypeTransaction
// 	sender := configuration.GetSelfServerAddress()
// 	messageHeader := types.NewMessageHeader(messageType, sender)
// 	messageBody := transactionData
// 	message := types.NewMessage(messageHeader, messageBody)
// 	return message
// }

// func sendArbitraryTransactionToNode(nodeConn types.NodeConnection) (bool, string) {
// 	conn := nodeConn.Conn
// 	message := createTransactionMessage()

// 	isMessageSent := comm.SendMessage(conn, message)

// 	return isMessageSent, message.Body.(string)
// }

// func SendArbitraryTransactionToAllNodes(connectionsList types.ConnectionsList) {
// 	for _, nodeConn := range connectionsList.GetNodeConnections() {
// 		isMessageSent, transactionData := sendArbitraryTransactionToNode(nodeConn)
// 		if !isMessageSent {
// 			log.Printf("Could not send arbitrary transaction %v to %v\n", transactionData, nodeConn.Node.GetAddress())
// 		} else {
// 			log.Printf("Sent arbitrary transaction %v to %v\n", transactionData, nodeConn.Node.GetAddress())
// 		}
// 	}
// }

func broadcastMessage(message types.Message, receivedFrom types.NodeAddress) {
	// Send to all the peers except the one from which the message was received
	_, connectionsList := configuration.ReadCurrentConnections("")
	for _, nodeConn := range connectionsList.GetNodeConnections() {
		if nodeConn.Node != receivedFrom {
			conn := nodeConn.Conn
			isMessageSent := comm.SendMessage(conn, message)
			if !isMessageSent {
				log.Printf("Could not send message %v to %v\n", message.Body.(string), nodeConn.Node.GetAddress())
			} else {
				log.Printf("Sent message %v to %v\n", message.Body.(string), nodeConn.Node.GetAddress())
			}
		}
	}
}

// Only following functions can be used by other packages

// Initiated by you as well as in response to messages received
func BroadcastTransaction(transaction types.Transaction, receivedFrom types.NodeAddress) {
	// Get self server address
	selfAddr := configuration.GetSelfServerAddress()
	message := comm.CreateMessage(selfAddr, types.MessageTypeTransaction, transaction)
	broadcastMessage(message, receivedFrom)
}

func BroadcastBlock(block types.Block, receivedFrom types.NodeAddress) {
	// Get self server address
	selfAddr := configuration.GetSelfServerAddress()
	message := comm.CreateMessage(selfAddr, types.MessageTypeBlock, block)
	broadcastMessage(message, receivedFrom)
}

// Only initiated by you
func BroadcastBlockRequest(blockHash string, receivedFrom types.NodeAddress) {
	// Get self server address
	selfAddr := configuration.GetSelfServerAddress()
	message := comm.CreateMessage(selfAddr, types.MessageTypeBlockRequest, blockHash)
	broadcastMessage(message, receivedFrom)
}

func SendBlockResponse(block types.Block, receivedFrom types.NodeAddress) {
	// Get self server address
	selfAddr := configuration.GetSelfServerAddress()
	message := comm.CreateMessage(selfAddr, types.MessageTypeBlockResponse, block)
	broadcastMessage(message, receivedFrom)
}

func BroadcastBlockChainRequest(receivedFrom types.NodeAddress) {
	// Get self server address
	selfAddr := configuration.GetSelfServerAddress()
	message := comm.CreateMessage(selfAddr, types.MessageTypeBlockChainRequest, "")
	broadcastMessage(message, receivedFrom)
}

func SendBlockChainResponse(blockChain types.BlockChain, receivedFrom types.NodeAddress) {
	// Get self server address
	selfAddr := configuration.GetSelfServerAddress()
	message := comm.CreateMessage(selfAddr, types.MessageTypeBlockChainResponse, blockChain)
	broadcastMessage(message, receivedFrom)
}
