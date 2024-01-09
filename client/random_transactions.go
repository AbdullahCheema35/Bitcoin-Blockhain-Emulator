package client

import (
	"fmt"
	"math/rand"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/nodestate"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/propagation"
)

func createRandomString() string {
	// Transaction body is randomly constructed string for now
	// TODO: Make a transaction type
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	length := 10
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	transactionStr := string(b)
	return transactionStr
}

func CreateRandomTransaction(counter int) {
	transactionStr := createRandomString()
	transactionStr = fmt.Sprintf("%d", counter) + "-" + transactionStr
	isAdded, newTx := nodestate.AddTransactionToPool(transactionStr)
	if isAdded {
		// Flood the transaction to all the peers except the one from which the transaction was received
		propagation.BroadcastTransaction(newTx, configuration.GetSelfServerAddress())
	}
}