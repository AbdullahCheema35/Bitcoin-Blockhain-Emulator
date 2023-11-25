package propagation

import "github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"

func CreateTransaction() {
	// Create a transaction

	// Take input from the user
	// Create a transaction from the input
	// Add the transaction to the transaction pool
	// Broadcast the transaction to all the peers
	BroadcastTransaction(transaction, selfNode)
}

func CreateBlock() {
	// Create a block from the transactions in the transaction pool
	// Max block size is N transactions
	// Min block size is M transactions
	// If the transaction pool has less than M transactions, wait for more transactions
	// If the transaction pool has more than N transactions, select N transactions from the pool
	// If the transaction pool has between M and N transactions, select all the transactions from the pool

	// Remove selected transactions from the transaction pool
	// Create a block from the selected transactions (using the function from part 1)

	// Start mining the block
	// When the block is mined, check if there is aleady a new block in the blockchain
	// If there is no block in the blockchain, add the mined block to the blockchain, and broadcast the block to all the peers
	// If there is a block in the blockchain, check if the mined block has a height greater than the height of the block in the blockchain
	// If the mined block has a height greater than the height of the block in the blockchain, add the mined block to the blockchain, and broadcast the block to all the peers
	// If the mined block has a height less than or equal to the height of the block in the blockchain, discard the mined block

	// If the mined block is discarded, add the transactions in the block to the transaction pool
	// If the mined block is added to the blockchain, remove the transactions in the block from the transaction pool

	// Broadcast the block to all the peers
	BroadcastBlock(block, selfNode)
}

func HandleReceivedTransaction(transaction types.Transaction, receivedFrom types.NodeAddress) {
	// Handle the received transaction from a node

	// Flood the transaction to all the peers except the one from which the transaction was received
	BroadcastTransaction(transaction, receivedFrom)
}

func HandleReceivedBlock(block types.Block, receivedFrom types.NodeAddress) {
	// Handle the received block from a node

	// Flood the block to all the peers except the one from which the block was received
	BroadcastBlock(block, receivedFrom)
}

func HandleBlockRequest(blockHash string, receivedFrom types.NodeAddress) {
	// Handle the received block request from a node

	// Send the requested block to the node from which the request was received
	SendBlockResponse(block, receivedFrom)
}

func HandleBlockChainRequest(receivedFrom types.NodeAddress) {
	// Handle the received blockchain request from a node

	// Send the blockchain to the node from which the request was received
	SendBlockChainResponse(blockChain, receivedFrom)
}
