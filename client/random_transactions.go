package client

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/nodestate"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/propagation"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
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
	selfPort := configuration.GetSelfServerAddress().Port
	transactionStr = fmt.Sprintf("%d--%d", counter, selfPort) + "-" + transactionStr
	isAdded, newTx := nodestate.AddTransactionToPool(transactionStr)
	if isAdded {
		// Flood the transaction to all the peers except the one from which the transaction was received
		propagation.BroadcastTransaction(newTx, configuration.GetSelfServerAddress())
	}
}

// CreateAndAddTransactionFromInput - Function to create and add a transaction from user input
func CreateTransaction() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter transaction value: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	// Trim newline character from the input
	value := strings.TrimSpace(input)

	// Create a new transaction using the NewTransaction function
	transaction := types.NewTransaction(value)

	// Assuming AddTransactionToPool is a function that adds a transaction to the pool
	// and returns a boolean status and the transaction object
	isAdded, transaction := nodestate.AddTransactionToPool(transaction.Value)

	if isAdded {
		// Flood the transaction to all the peers except the one from which the transaction was received
		propagation.BroadcastTransaction(transaction, configuration.GetSelfServerAddress())
	}
}
