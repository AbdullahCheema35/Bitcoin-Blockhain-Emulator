package types

// BlockNode struct representing a node in the blockchain
type BlockNode struct {
	Block    Block
	PrevNode *BlockNode
}

// BlockChain struct representing the blockchain
type BlockChain struct {
	LatestNode *BlockNode
}

// NewBlockChain constructs a blockchain object
func NewBlockChain() BlockChain {
	return BlockChain{
		LatestNode: nil,
	}
}

func (bchain *BlockChain) Display() {
	currentNode := bchain.LatestNode
	for currentNode != nil {
		currentNode.Block.Display()
		currentNode = currentNode.PrevNode
	}
}

// AddBlock adds a block to the blockchain, if the height of the block is greater than the height of the latest block in the blockchain
func (bchain *BlockChain) AddBlock(b Block) bool {
	if bchain.LatestNode == nil {
		bchain.LatestNode = &BlockNode{
			Block:    b,
			PrevNode: nil,
		}
		return true
	} else if bchain.LatestNode.Block.Header.Height < b.Header.Height {
		if bchain.LatestNode.Block.Header.Height == b.Header.Height-1 {
			if bchain.LatestNode.Block.BlockHash == b.Header.PreviousBlockHash {
				// Insert the block as the latest node
				newBlockNode := &BlockNode{
					Block:    b,
					PrevNode: bchain.LatestNode,
				}
				bchain.LatestNode = newBlockNode
				return true
			} else {
				return false // Previous block hash doesn't match
			}
		} else {
			// Insert the block as the latest node
			newBlockNode := &BlockNode{
				Block:    b,
				PrevNode: bchain.LatestNode,
			}
			bchain.LatestNode = newBlockNode
			return true
		}
	}
	return false
}

// func CreateNilBlock(nilBlockHash string) Block {
// 	blockHeader := BlockHeader{
// 		PreviousBlockHash: nilBlockHash,
// 		Nonce:             -1,
// 		Height:            -1,
// 		MerkleRootHash:    "",
// 		Timestamp:         time.Now(),
// 		DifficultyTarget:  -1,
// 	}
// 	blockBody := BlockBody{
// 		Transactions: NewTransactionList(),
// 		MerkleTree:   MerkleTree{},
// 	}
// 	return Block{
// 		BlockHash: "",
// 		Header:    blockHeader,
// 		Body:      blockBody,
// 	}
// }

func (bchain *BlockChain) GetLatestBlockNode() *BlockNode {
	// if bchain.LatestNode == nil {
	// 	return CreateNilBlock(bchain.GetLatestBlockHash())
	// } else {
	// 	return bchain.LatestNode.Block
	// }
	return bchain.LatestNode
}

func (bchain *BlockChain) GetLatestBlockHeight() int {
	if bchain.LatestNode == nil {
		return -1
	} else {
		return bchain.LatestNode.Block.Header.Height
	}
}

func (bchain *BlockChain) GetLatestBlockHash() string {
	if bchain.LatestNode == nil {
		return "0000000000000000000000000000000000000000000000000000000000000000"
	} else {
		return bchain.LatestNode.Block.BlockHash
	}
}

// // RemoveBlock removes a block from the blockchain based on its height
// func (bchain *BlockChain) RemoveBlock(height int) bool {
// 	if bchain.LatestNode == nil {
// 		return false // Blockchain is empty, no blocks to remove
// 	}

// 	// If the block to remove is the latest block
// 	if bchain.LatestNode.Block.Header.Height == height {
// 		bchain.LatestNode = bchain.LatestNode.PrevNode
// 		return true
// 	}

// 	// Traverse through the blockchain to find the block with the specified height
// 	prevNode := bchain.LatestNode
// 	currentNode := bchain.LatestNode.PrevNode

// 	for currentNode != nil && currentNode.Block.Header.Height >= height {
// 		if currentNode.Block.Header.Height == height {
// 			prevNode.PrevNode = currentNode.PrevNode
// 			return true
// 		}
// 		prevNode = currentNode
// 		currentNode = currentNode.PrevNode
// 	}

// 	return false // Block with specified height not found
// }

// // UpdateBlockChain adds a missing block to the blockchain (Blockchain mustn't be empty)
// // It will only add the block if the block height is less than the height of the latest block in the blockchain
// // It will only add the block if the block doesn't already exist in the blockchain
// // It will only add the block if hashes of adjacent blocks match
// func (bchain *BlockChain) AddMissingBlock(b Block) bool {
// 	if b.Header.Height < 0 {
// 		return false // Invalid block height
// 	} else if bchain.LatestNode == nil {
// 		return false // Blockchain is empty
// 	} else if b.Header.Height >= bchain.LatestNode.Block.Header.Height {
// 		return false // Block height is greater than or equal to the height of the latest block in the blockchain
// 	}

// 	if bchain.SameBlockExists(b) {
// 		return false // Same block already exists at the same height in the blockchain
// 	}

// 	// Find the appropriate position based on height and ensure adjacent hashes match
// 	prevNode := bchain.LatestNode
// 	currentNode := bchain.LatestNode

// 	for currentNode != nil && currentNode.Block.Header.Height > b.Header.Height {
// 		prevNode = currentNode
// 		currentNode = currentNode.PrevNode
// 	}

// 	if currentNode == nil {
// 		if prevNode.Block.Header.PreviousBlockHash == b.BlockHash {
// 			// Insert the block as the last node
// 			newBlockNode := &BlockNode{
// 				Block:    b,
// 				PrevNode: nil,
// 			}
// 			prevNode.PrevNode = newBlockNode
// 			return true
// 		}
// 	} else if currentNode.Block.Header.Height < b.Header.Height {
// 		// Insert the block in between two nodes
// 		newBlockNode := &BlockNode{
// 			Block:    b,
// 			PrevNode: currentNode,
// 		}
// 		prevNode.PrevNode = newBlockNode
// 		return true
// 	}
// 	return false
// }

// Helper function to check if same block already exists in the blockchain
func (bchain *BlockChain) SameBlockExists(b Block) bool {
	currentNode := bchain.LatestNode

	for currentNode != nil && currentNode.Block.Header.Height >= b.Header.Height {
		if currentNode.Block.Header.Height == b.Header.Height && currentNode.Block.BlockHash == b.BlockHash {
			return true
		}
		currentNode = currentNode.PrevNode
	}
	return false
}

// Helper function to check if a block with the same height already exists in the blockchain
func (bchain *BlockChain) BlockExistsAtHeight(b Block) bool {
	currentNode := bchain.LatestNode

	for currentNode != nil && currentNode.Block.Header.Height >= b.Header.Height {
		if currentNode.Block.Header.Height == b.Header.Height {
			return true
		}
		currentNode = currentNode.PrevNode
	}
	return false
}
