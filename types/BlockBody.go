package types

// BlockBody struct representing the body of a block
type BlockBody struct {
	Transactions TransactionList
	MerkleTree   MerkleTree // Contains the merkle tree of the transactions
}

type MerkleTree struct {
	RootHash string
}
