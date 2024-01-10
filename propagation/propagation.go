package propagation

import (
	"log"
	"time"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/comm"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/nodestate"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

func broadcastMessage(message types.Message, receivedFrom types.NodeAddress) {
	_, connectionsList := nodestate.ReadCurrentConnections("")
	for _, nodeConn := range connectionsList.GetNodeConnections() {
		if nodeConn.Node != receivedFrom {
			isMessageSent := comm.SendMessage(nodeConn, message)
			if !isMessageSent {
				log.Printf("Could not send message of type %v to %v\n", message.Header.Type, nodeConn.Node.GetAddress())
			} else {
				log.Printf("Sent message %v to %v\n", message.Body, nodeConn.Node.GetAddress())
			}
		}
	}
}

func sendResponse(message types.Message, receivedFrom types.NodeAddress) {
	_, connectionsList := nodestate.ReadCurrentConnections("")
	for _, nodeConn := range connectionsList.GetNodeConnections() {
		if nodeConn.Node == receivedFrom {
			isMessageSent := comm.SendMessage(nodeConn, message)
			if !isMessageSent {
				log.Printf("Could not send response %v to %v\n", message.Body, nodeConn.Node.GetAddress())
			} else {
				log.Printf("Sent response %v to %v\n", message.Body, nodeConn.Node.GetAddress())
			}
			return
		}
	}
}

func BroadcastTransaction(transaction types.Transaction, receivedFrom types.NodeAddress) {
	selfAddr := configuration.GetSelfServerAddress()
	message := comm.CreateMessage(selfAddr, types.MessageTypeTransaction, transaction)
	broadcastMessage(message, receivedFrom)
}

func BroadcastBlock(block types.Block, receivedFrom types.NodeAddress) {
	selfAddr := configuration.GetSelfServerAddress()
	message := comm.CreateMessage(selfAddr, types.MessageTypeBlock, block)
	broadcastMessage(message, receivedFrom)
}

func BroadcastBlockChainRequest() {
	selfAddr := configuration.GetSelfServerAddress()
	message := comm.CreateMessage(selfAddr, types.MessageTypeBlockChainRequest, nil)
	broadcastMessage(message, selfAddr)
}

func SendBlockChainResponse(blockChain types.BlockChain, receivedFrom types.NodeAddress) {
	selfAddr := configuration.GetSelfServerAddress()
	message := comm.CreateMessage(selfAddr, types.MessageTypeBlockChainResponse, blockChain)
	sendResponse(message, receivedFrom)
}

func InitiateTopologyRequest() types.TopologyRequest {
	selfAddr := configuration.GetSelfServerAddress()

	nodesFound := types.NewNodesList()
	nodesFound.AddNode(selfAddr)

	_, currentConns := nodestate.ReadCurrentConnections("")
	peersList := types.NewNodesList()

	for _, nodeConn := range currentConns.GetNodeConnections() {
		nodesFound.AddNode(nodeConn.Node)
		peersList.AddNode(nodeConn.Node)
	}

	networkList := types.NewNetworkList(selfAddr, peersList)
	originList := types.NewNodesList()
	originList.AddNode(selfAddr)

	topologyRequest := types.NewTopologyRequest(originList, nodesFound, networkList)
	message := comm.CreateMessage(selfAddr, types.MessageTypeTopologyRequest, topologyRequest)
	broadcastMessage(message, selfAddr)

	return topologyRequest
}

func BroadcastTopologyRequest(topologyRequest types.TopologyRequest, toSendPeersList types.NodesList) {
	selfAddr := configuration.GetSelfServerAddress()
	message := comm.CreateMessage(selfAddr, types.MessageTypeTopologyRequest, topologyRequest)

	_, connectionsList := nodestate.ReadCurrentConnections("")
	for _, nodeConn := range connectionsList.GetNodeConnections() {
		if toSendPeersList.ContainsNode(nodeConn.Node) {
			isMessageSent := comm.SendMessage(nodeConn, message)
			if !isMessageSent {
				log.Printf("Could not send topology request to %s\n", nodeConn.Node.GetAddress())
			} else {
				log.Printf("Sent topology request to %s\n", nodeConn.Node.GetAddress())
			}
		}
	}
}

func SendTopologyResponse(topologyRequest types.TopologyRequest, receivedFrom types.NodeAddress) {
	selfAddr := configuration.GetSelfServerAddress()
	message := comm.CreateMessage(selfAddr, types.MessageTypeTopologyResponse, topologyRequest)

	_, connectionsList := nodestate.ReadCurrentConnections("")
	for _, nodeConn := range connectionsList.GetNodeConnections() {
		if nodeConn.Node == receivedFrom {
			isMessageSent := comm.SendMessage(nodeConn, message)
			if !isMessageSent {
				log.Printf("Could not send topology response to %s\n", nodeConn.Node.GetAddress())
			} else {
				log.Printf("Sent topology response to %s\n", nodeConn.Node.GetAddress())
			}
			return
		}
	}

	log.Println("Couldn't send topology response received from", receivedFrom.GetAddress())
}

func GetP2PNetwork() []types.NetworkList {
	topologyChan := nodestate.InitTopologyChan()
	topologyRequest := InitiateTopologyRequest()

	nodesFound := topologyRequest.NodesFound
	networkList := topologyRequest.ThisNodePeers

	listOfNetworkList := make([]types.NetworkList, 0)
	listOfNetworkList = append(listOfNetworkList, networkList)

	for len(nodesFound.GetNodes()) != len(listOfNetworkList) {
		select {
		case topologyRequest := <-topologyChan:
			newNodesFound := topologyRequest.NodesFound
			for _, node := range newNodesFound.GetNodes() {
				nodesFound.AddNode(node)
			}
			networkList := topologyRequest.ThisNodePeers
			listOfNetworkList = append(listOfNetworkList, networkList)
		default:
			time.Sleep(500 * time.Millisecond)
		}
	}

	nodestate.CloseTopologyChan()
	return listOfNetworkList
}
