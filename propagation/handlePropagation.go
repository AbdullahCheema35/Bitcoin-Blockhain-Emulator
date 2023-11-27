package propagation

import (
	"log"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/mining"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/nodestate"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

// import "github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"

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
	// Check if the transaction is already present in the transaction pool
	// If the transaction is not present in the transaction pool, add the transaction to the transaction pool
	// If the transaction is already present in the transaction pool, discard the transaction

	transactionPool := nodestate.LockTransactionPool()
	isAdded := transactionPool.AddTransaction(transaction)
	nodestate.UnlockTransactionPool(transactionPool)

	if isAdded {
		// Flood the transaction to all the peers except the one from which the transaction was received
		BroadcastTransaction(transaction, receivedFrom)
	}
	// else discard the transaction
}

func HandleReceivedBlock(block types.Block, receivedFrom types.NodeAddress) {
	// Handle the received block from a node
	// Verify the received block
	// If the block is valid, add the block to the blockchain, and broadcast the block to all the peers
	// If the block is invalid, discard the block

	result := mining.AddNewBlockToBlockChain(block)

	switch result {
	case mining.NewBlockVerificationFailed:
		return
	case mining.NewHeightLEQCurrentHeight:
		return
	case mining.NewBlockDuplicateTransactions:
		selfNode := configuration.GetSelfServerAddress()
		if selfNode.GetAddress() == receivedFrom.GetAddress() {
			return
		}
		BroadcastBlockChainRequest()
		return
	case mining.NewBlockPrevHashDontMatch:
		selfNode := configuration.GetSelfServerAddress()
		if selfNode.GetAddress() == receivedFrom.GetAddress() {
			return
		}
		BroadcastBlockChainRequest()
		return
	case mining.NewBlockAddedSuccessfully:
		go StartCreateBlock()
		// Flood the block to all the peers except the one from which the block was received
		BroadcastBlock(block, receivedFrom)
	}
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

func HandleTopologyRequest(receivedFrom types.NodeAddress, topologyRequest types.TopologyRequest) {
	// Handle the received topology request from a node
	// Add my own peers to the list of peers in the topology request
	// Broadcast the topology request to all the peers except the one from which the request was received and the ones already present in the topology request

	// Get self server address
	selfAddr := configuration.GetSelfServerAddress()

	myPeersList := types.NewNodesList()
	toSendPeersList := types.NewNodesList()

	previousPeersList := topologyRequest.NodesFound

	_, currentConns := nodestate.ReadCurrentConnections("")

	// Get this node's peers and add them to the list
	for _, nodeConn := range currentConns.GetNodeConnections() {
		myPeersList.AddNode(nodeConn.Node)
		added := previousPeersList.AddNode(nodeConn.Node)
		if added {
			toSendPeersList.AddNode(nodeConn.Node)
		}
	}

	networkList := types.NewNetworkList(selfAddr, myPeersList)

	// Add current node to the origin list
	originList := topologyRequest.Origin
	originList.AddNode(selfAddr)

	// Create topology request object
	newTopologyRequest := types.NewTopologyRequest(originList, previousPeersList, networkList)

	SendTopologyResponse(newTopologyRequest, receivedFrom)
	BroadcastTopologyRequest(newTopologyRequest, toSendPeersList)
}

func HandleTopologyResponse(topologyRequest types.TopologyRequest, receivedFrom types.NodeAddress) {
	// Handle the received topology response from a node
	// Add the node from which the response was received to the list of nodes found in the topology request
	// Add the peers of the node from which the response was received to the list of nodes found in the topology request
	// Broadcast the topology request to all the peers except the one from which the request was received and the ones already present in the topology request

	// Get self server address
	selfAddr := configuration.GetSelfServerAddress()

	originList := topologyRequest.Origin
	// Pop the last node from the origin list
	originList.Nodes = originList.Nodes[:len(originList.Nodes)-1]

	// Check if I am the origin of the topology request (i.e. the topology request has completed its journey)
	// Check if length of origin list is 1
	if len(originList.Nodes) != 1 {
		// Send the topology response to the next node in the origin list
		nextNode := originList.Nodes[len(originList.Nodes)-1]
		SendTopologyResponse(topologyRequest, nextNode)
		return
	}
	// Check if the origin of the topology request is the same as the node from which the response was received
	if originList.Nodes[0].GetAddress() != selfAddr.GetAddress() {
		log.Println("\n\n\n\n\n\nLost my way\n\n\n\n\n\n")
		return
	}

	// The topology request has completed its journey
	topologyChan := nodestate.GetTopologyChan()

	select {
	case val, ok := <-topologyChan:
		if !ok {
			log.Println("Topology channel closed")
			return
		} else {
			topologyChan <- val
		}
	default:
		topologyChan <- topologyRequest
	}
}
