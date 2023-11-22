package propagation

import (
	"log"
	"math/rand"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/common"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

func createTransactionMessage() types.Message {
	// Transaction body is randomly constructed string for now
	// TODO: Make a transaction type
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	length := 10
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	transactionData := string(b)
	// End of random string generation
	// Temp fix

	messageType := types.MessageTypeTransaction
	sender := configuration.GetSelfServerAddress()
	messageHeader := types.NewMessageHeader(messageType, sender)
	messageBody := transactionData
	message := types.NewMessage(messageHeader, messageBody)
	return message
}

func sendArbitraryTransactionToNode(nodeConn types.NodeConnection) (bool, string) {
	conn := nodeConn.Conn
	message := createTransactionMessage()

	isMessageSent := common.SendMessage(conn, message)

	return isMessageSent, message.Body.(string)
}

func SendArbitraryTransactionToAllNodes(connectionsList types.ConnectionsList) {
	for _, nodeConn := range connectionsList.GetNodeConnections() {
		isMessageSent, transactionData := sendArbitraryTransactionToNode(nodeConn)
		if !isMessageSent {
			log.Printf("Could not send arbitrary transaction %v to %v\n", transactionData, nodeConn.Node.GetAddress())
		} else {
			log.Printf("Sent arbitrary transaction %v to %v\n", transactionData, nodeConn.Node.GetAddress())
		}
	}
}
