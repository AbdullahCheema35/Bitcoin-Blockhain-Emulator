package types

import (
	"crypto/sha256"
	"encoding/hex"
)

type Transaction struct {
	Value string
	Hash  string
}

// Function to calculate the hash of a transaction value
func (tx *Transaction) RecalculateTransactionHash() string {
	hash := sha256.Sum256([]byte(tx.Value))
	hashString := hex.EncodeToString(hash[:])
	tx.Hash = hashString
	return hashString
}

func NewTransaction(value string) Transaction {
	transaction := Transaction{
		Value: value,
	}
	transaction.RecalculateTransactionHash()
	return transaction
}

func (tx *Transaction) ModifyTransaction(value string) {
	tx.Value = value
	tx.RecalculateTransactionHash()
}
