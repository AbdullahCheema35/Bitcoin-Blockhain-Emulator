package types

import (
	"fmt"
)

type TransactionList struct {
	Transactions []Transaction
}

func NewTransactionList() TransactionList {
	return TransactionList{Transactions: make([]Transaction, 0)}
}

func (tl *TransactionList) AddTransaction(transaction Transaction) bool {
	for _, t := range tl.Transactions {
		if t.Hash == transaction.Hash {
			return false
		}
	}
	tl.Transactions = append(tl.Transactions, transaction)
	return true
}

func (tl *TransactionList) RemoveTransaction(transaction Transaction) bool {
	for i, t := range tl.Transactions {
		if t.Hash == transaction.Hash {
			tl.Transactions = append(tl.Transactions[:i], tl.Transactions[i+1:]...)
			return true
		}
	}
	return false
}

func (tl *TransactionList) GetTransactions() []Transaction {
	return tl.Transactions
}

// Check if the transaction list contains a transaction
func (tl *TransactionList) Contains(transaction Transaction) bool {
	for _, t := range tl.Transactions {
		if t.Hash == transaction.Hash {
			return true
		}
	}
	return false
}

func DisplayTransactionPool(transactionPool *TransactionList) {
	fmt.Println("Transaction Pool:")
	transactions := transactionPool.GetTransactions()

	if len(transactions) == 0 {
		fmt.Println("No transactions in the pool.")
		return
	}

	for i, transaction := range transactions {
		fmt.Printf("%d. Hash: %s, Value: %s\n", i+1, transaction.Hash, transaction.Value)
	}
	fmt.Println("------------------------")
}

func NewTransactionListFromSlice(transactions []Transaction) TransactionList {
	return TransactionList{Transactions: transactions}
}
