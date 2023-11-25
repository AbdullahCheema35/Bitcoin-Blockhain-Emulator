package types

// BlockChain struct representing the blockchain
type BlockChain struct {
	LastNode *BlockNode
}

// BlockNode struct representing a node in the blockchain
type BlockNode struct {
	Block    Block
	PrevNode *BlockNode
}
