package propagation

import (
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/mining"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/nodestate"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

func HandleReceivedTransaction(transaction types.Transaction, receivedFrom types.NodeAddress) {
	isAdded, addedTx := nodestate.AddTransactionToPool(transaction.Value)

	if isAdded {
		// log.Println("Transaction added to the pool, broadcasting to peers.")
		BroadcastTransaction(addedTx, receivedFrom)
	}
}

func HandleReceivedBlock(block types.Block, receivedFrom types.NodeAddress) {
	returnType := mining.HandleNewBlock(block, receivedFrom)
	switch returnType {
	case types.DoNothing:
		// log.Println("No action required for the received block.")
	case types.InitiateBroadcastBlockChainRequest:
		// log.Println("Initiating broadcast of the block chain request.")
		BroadcastBlockChainRequest()
	case types.InitiateBroadcastBlock:
		// log.Println("Initiating broadcast of the received block.")
		BroadcastBlock(block, receivedFrom)
	}
	// // Display all the blocks
	// bchain := nodestate.ReadBlockChain()
	// bchain.DisplayHeaderInfo()

}

func HandleReceivedBlockChain(blockChain types.BlockChain, receivedFrom types.NodeAddress) {
	mining.HandleNewBChain(blockChain)
}

func HandleBlockChainRequest(receivedFrom types.NodeAddress) {
	bchain := nodestate.ReadBlockChain()
	SendBlockChainResponse(bchain, receivedFrom)
}

func HandleTopologyRequest(receivedFrom types.NodeAddress, topologyRequest types.TopologyRequest) {
	selfAddr := configuration.GetSelfServerAddress()
	myPeersList := types.NewNodesList()
	toSendPeersList := types.NewNodesList()

	previousNodesList := topologyRequest.NodesFound
	_, currentConns := nodestate.ReadCurrentConnections("")

	for _, nodeConn := range currentConns.GetNodeConnections() {
		myPeersList.AddNode(nodeConn.Node)
		added := previousNodesList.AddNode(nodeConn.Node)
		if added {
			toSendPeersList.AddNode(nodeConn.Node)
		}
	}

	networkList := types.NewNetworkList(selfAddr, myPeersList)
	originList := topologyRequest.Origin
	originList.AddNode(selfAddr)

	newTopologyRequest := types.NewTopologyRequest(originList, previousNodesList, networkList)
	SendTopologyResponse(newTopologyRequest, receivedFrom)
	BroadcastTopologyRequest(newTopologyRequest, toSendPeersList)
}

func HandleTopologyResponse(topologyRequest types.TopologyRequest, receivedFrom types.NodeAddress) {
	selfAddr := configuration.GetSelfServerAddress()
	originList := topologyRequest.Origin

	for len(originList.Nodes) > 0 && originList.Nodes[len(originList.Nodes)-1].GetAddress() != selfAddr.GetAddress() {
		originList.Nodes = originList.Nodes[:len(originList.Nodes)-1]
	}

	if len(originList.Nodes) == 0 {
		return
	}

	if len(originList.Nodes) != 1 {
		originList.Nodes = originList.Nodes[:len(originList.Nodes)-1]
		nextNode := originList.Nodes[len(originList.Nodes)-1]
		SendTopologyResponse(topologyRequest, nextNode)
		return
	}

	topologyChan := nodestate.GetTopologyChan()

	select {
	case val, ok := <-topologyChan:
		if !ok {
			return
		} else {
			topologyChan <- val
			topologyChan <- topologyRequest
		}
	default:
		topologyChan <- topologyRequest
	}
}
