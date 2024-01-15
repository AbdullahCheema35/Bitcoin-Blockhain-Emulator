package propagation

import (
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/mining"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/nodestate"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

func HandleReceivedTransaction(transaction types.Transaction, receivedFrom types.NodeAddress) {
	// Handle the received transaction from a node
	// Check if the transaction is already present in the transaction pool
	// If the transaction is not present in the transaction pool, add the transaction to the transaction pool
	// If the transaction is already present in the transaction pool, discard the transaction

	isAdded, addedTx := nodestate.AddTransactionToPool(transaction.Value)

	if isAdded {
		// Flood the transaction to all the peers except the one from which the transaction was received
		BroadcastTransaction(addedTx, receivedFrom)
	}
	// else discard the transaction
}

func HandleReceivedBlock(block types.Block, receivedFrom types.NodeAddress) {
	returnType := mining.HandleNewBlock(block, receivedFrom)
	switch returnType {
	case types.DoNothing:
		// fmt.Println("Do nothing")
	case types.InitiateBroadcastBlockChainRequest:
		// fmt.Println("InitiateBroadcastBlockChainRequest")
		BroadcastBlockChainRequest()
	case types.InitiateBroadcastBlock:
		// fmt.Println("InitiateBroadcastBlock")
		BroadcastBlock(block, receivedFrom)
	}

	// Display all the blocks
	bchain := nodestate.ReadBlockChain()
	bchain.Display()
}

func HandleReceivedBlockChain(blockChain types.BlockChain, receivedFrom types.NodeAddress) {
	// Handle the received blockchain from a node
	// Verify the received blockchain
	// If the blockchain is valid, replace the current blockchain with the received blockchain
	// If the blockchain is invalid, discard the blockchain
	mining.HandleNewBChain(blockChain)
}

// func HandleBlockRequest(blockHash string, receivedFrom types.NodeAddress) {
// 	// Handle the received block request from a node

// 	// Send the requested block to the node from which the request was received
// 	SendBlockResponse(block, receivedFrom)
// }

func HandleBlockChainRequest(receivedFrom types.NodeAddress) {
	// Handle the received blockchain request from a node

	// Get the current blockchain
	bchain := nodestate.ReadBlockChain()

	// Send the blockchain to the node from which the request was received
	SendBlockChainResponse(bchain, receivedFrom)
}

func HandleTopologyRequest(receivedFrom types.NodeAddress, topologyRequest types.TopologyRequest) {
	// Handle the received topology request from a node
	// Add my own peers to the list of peers in the topology request
	// Broadcast the topology request to all the peers except the one from which the request was received and the ones already present in the topology request

	// Get self server address
	selfAddr := configuration.GetSelfServerAddress()

	myPeersList := types.NewNodesList()
	toSendPeersList := types.NewNodesList()

	previousNodesList := topologyRequest.NodesFound

	_, currentConns := nodestate.ReadCurrentConnections("")

	// Get this node's peers and add them to the list
	for _, nodeConn := range currentConns.GetNodeConnections() {
		myPeersList.AddNode(nodeConn.Node)
		added := previousNodesList.AddNode(nodeConn.Node)
		if added {
			toSendPeersList.AddNode(nodeConn.Node)
		}
	}

	networkList := types.NewNetworkList(selfAddr, myPeersList)

	// Add current node to the origin list
	originList := topologyRequest.Origin
	originList.AddNode(selfAddr)

	// Create topology request object
	newTopologyRequest := types.NewTopologyRequest(originList, previousNodesList, networkList)
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

	// OriginList
	originList := topologyRequest.Origin

	// keep popping from the list until selfAddr found or list empty
	for len(originList.Nodes) > 0 && originList.Nodes[len(originList.Nodes)-1].GetAddress() != selfAddr.GetAddress() {
		originList.Nodes = originList.Nodes[:len(originList.Nodes)-1]
	}

	// If the origin list is empty, then the topology request has lost its way
	if len(originList.Nodes) == 0 {
		// log.Println("\n\n\n\n\n\nLost my way\n\n\n\n\n\n")
		return
	}

	// Last element is the same as me (same address)

	// Check if I am the origin of the topology request (i.e. the topology request has completed its journey)
	// Check if length of origin list is not 1 (Im not origin)
	if len(originList.Nodes) != 1 {
		// Pop the last node from the origin list, which is the same node as this node
		originList.Nodes = originList.Nodes[:len(originList.Nodes)-1]
		// Send the topology request to the next node in the origin list
		nextNode := originList.Nodes[len(originList.Nodes)-1]
		// Send the node
		SendTopologyResponse(topologyRequest, nextNode)
		return
	}

	// // Pop the last node from the origin list
	// originList.Nodes = originList.Nodes[:len(originList.Nodes)-1]

	// log.Println("Topology request has completed its journey")

	// The topology request has completed its journey
	topologyChan := nodestate.GetTopologyChan()

	// Send the topology request to the topology channel
	// // // log.Println("Handle Topology response")
	// topologyChan <- topologyRequest

	select {
	case val, ok := <-topologyChan:
		if !ok {
			// log.Println("Handle Topology response: case !ok")
			// log.Println("Topology channel closed")
			return
		} else {
			// log.Println("Handle Topology response: case ok")
			topologyChan <- val
			topologyChan <- topologyRequest
		}
	default:
		// log.Println("Handle Topology response: default")
		topologyChan <- topologyRequest
	}
}
