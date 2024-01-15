package propagation

import (
	"time"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/comm"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/configuration"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/nodestate"
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

// func createTransactionMessage() types.Message {
// 	// Transaction body is randomly constructed string for now
// 	// TODO: Make a transaction type
// 	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
// 	length := 10
// 	b := make([]rune, length)
// 	for i := range b {
// 		b[i] = letters[rand.Intn(len(letters))]
// 	}
// 	transactionData := string(b)
// 	// End of random string generation
// 	// Temp fix

// 	messageType := types.MessageTypeTransaction
// 	sender := configuration.GetSelfServerAddress()
// 	messageHeader := types.NewMessageHeader(messageType, sender)
// 	messageBody := transactionData
// 	message := types.NewMessage(messageHeader, messageBody)
// 	return message
// }

// func sendArbitraryTransactionToNode(nodeConn types.NodeConnection) (bool, string) {
// 	conn := nodeConn.Conn
// 	message := createTransactionMessage()

// 	isMessageSent := comm.SendMessage(conn, message)

// 	return isMessageSent, message.Body.(string)
// }

// func SendArbitraryTransactionToAllNodes(connectionsList types.ConnectionsList) {
// 	for _, nodeConn := range connectionsList.GetNodeConnections() {
// 		isMessageSent, transactionData := sendArbitraryTransactionToNode(nodeConn)
// 		if !isMessageSent {
// 			// log.Printf("Could not send arbitrary transaction %v to %v\n", transactionData, nodeConn.Node.GetAddress())
// 		} else {
// 			// log.Printf("Sent arbitrary transaction %v to %v\n", transactionData, nodeConn.Node.GetAddress())
// 		}
// 	}
// }

func broadcastMessage(message types.Message, receivedFrom types.NodeAddress) {
	// Send to all the peers except the one from which the message was received
	_, connectionsList := nodestate.ReadCurrentConnections("")
	for _, nodeConn := range connectionsList.GetNodeConnections() {
		if nodeConn.Node != receivedFrom {
			conn := nodeConn.Conn
			isMessageSent := comm.SendMessage(conn, message)
			if !isMessageSent {
				// log.Printf("Could not send message %v to %v\n", message.Body, nodeConn.Node.GetAddress())
			} else {
				// log.Printf("Sent message %v to %v\n", message.Body, nodeConn.Node.GetAddress())
			}
		}
	}
}

func sendResponse(message types.Message, receivedFrom types.NodeAddress) {
	// send to the node from which the request was received
	_, connectionsList := nodestate.ReadCurrentConnections("")
	for _, nodeConn := range connectionsList.GetNodeConnections() {
		if nodeConn.Node == receivedFrom {
			conn := nodeConn.Conn
			isMessageSent := comm.SendMessage(conn, message)
			if !isMessageSent {
				// log.Printf("Could not send response %v to %v\n", message.Body, nodeConn.Node.GetAddress())
			} else {
				// log.Printf("Sent response %v to %v\n", message.Body, nodeConn.Node.GetAddress())
			}
			return
		}
	}
}

// Only following functions can be used by other packages

// Initiated by you as well as in response to messages received
func BroadcastTransaction(transaction types.Transaction, receivedFrom types.NodeAddress) {
	// Get self server address
	selfAddr := configuration.GetSelfServerAddress()
	message := comm.CreateMessage(selfAddr, types.MessageTypeTransaction, transaction)
	broadcastMessage(message, receivedFrom)
}

func BroadcastBlock(block types.Block, receivedFrom types.NodeAddress) {
	// Get self server address
	selfAddr := configuration.GetSelfServerAddress()
	message := comm.CreateMessage(selfAddr, types.MessageTypeBlock, block)
	broadcastMessage(message, receivedFrom)
}

// // Only initiated by you
// func BroadcastBlockRequest(blockHash string, receivedFrom types.NodeAddress) {
// 	// Get self server address
// 	selfAddr := configuration.GetSelfServerAddress()
// 	message := comm.CreateMessage(selfAddr, types.MessageTypeBlockRequest, blockHash)
// 	broadcastMessage(message, receivedFrom)
// }

// func SendBlockResponse(block types.Block, receivedFrom types.NodeAddress) {
// 	// Get self server address
// 	selfAddr := configuration.GetSelfServerAddress()
// 	message := comm.CreateMessage(selfAddr, types.MessageTypeBlockResponse, block)
// 	broadcastMessage(message, receivedFrom)
// }

func BroadcastBlockChainRequest() {
	// Get self server address
	selfAddr := configuration.GetSelfServerAddress()
	message := comm.CreateMessage(selfAddr, types.MessageTypeBlockChainRequest, nil)
	broadcastMessage(message, selfAddr)
}

func SendBlockChainResponse(blockChain types.BlockChain, receivedFrom types.NodeAddress) {
	// Get self server address
	selfAddr := configuration.GetSelfServerAddress()
	message := comm.CreateMessage(selfAddr, types.MessageTypeBlockChainResponse, blockChain)
	sendResponse(message, receivedFrom)
}

func InitiateTopologyRequest() types.TopologyRequest {
	// Get self server address
	selfAddr := configuration.GetSelfServerAddress()

	nodesFound := types.NewNodesList()
	nodesFound.AddNode(selfAddr)

	_, currentConns := nodestate.ReadCurrentConnections("")

	peersList := types.NewNodesList()

	// Get this node's peers and add them to the list
	for _, nodeConn := range currentConns.GetNodeConnections() {
		nodesFound.AddNode(nodeConn.Node)
		peersList.AddNode(nodeConn.Node)
	}

	networkList := types.NewNetworkList(selfAddr, peersList)

	originList := types.NewNodesList()
	originList.AddNode(selfAddr)

	// Create topology request object
	topologyRequest := types.NewTopologyRequest(originList, nodesFound, networkList)

	message := comm.CreateMessage(selfAddr, types.MessageTypeTopologyRequest, topologyRequest)
	broadcastMessage(message, selfAddr)

	return topologyRequest
}

func BroadcastTopologyRequest(topologyRequest types.TopologyRequest, toSendPeersList types.NodesList) {
	// Get self server address
	selfAddr := configuration.GetSelfServerAddress()

	message := comm.CreateMessage(selfAddr, types.MessageTypeTopologyRequest, topologyRequest)

	// Send to all the peers except the ones already present in the topology request
	_, connectionsList := nodestate.ReadCurrentConnections("")
	for _, nodeConn := range connectionsList.GetNodeConnections() {
		if toSendPeersList.ContainsNode(nodeConn.Node) {
			conn := nodeConn.Conn
			isMessageSent := comm.SendMessage(conn, message)
			if !isMessageSent {
				// log.Printf("Line 162: Could not send topology request to %s\n", nodeConn.Node.GetAddress())
			} else {
				// log.Printf("Line 164: Sent topology request to %s\n", nodeConn.Node.GetAddress())
			}
		}
	}
}

func SendTopologyResponse(topologyRequest types.TopologyRequest, receivedFrom types.NodeAddress) {
	// Get self server address
	selfAddr := configuration.GetSelfServerAddress()

	message := comm.CreateMessage(selfAddr, types.MessageTypeTopologyResponse, topologyRequest)

	// Send to the node from which the request was received
	_, connectionsList := nodestate.ReadCurrentConnections("")
	for _, nodeConn := range connectionsList.GetNodeConnections() {
		if nodeConn.Node == receivedFrom {
			conn := nodeConn.Conn
			isMessageSent := comm.SendMessage(conn, message)
			if !isMessageSent {
				// log.Printf("Line 183: Could not send topology response to %s\n", nodeConn.Node.GetAddress())
			} else {
				// log.Printf("Line 185: Sent topology response to %s\n", nodeConn.Node.GetAddress())
			}
			return
		}
	}

	// log.Println("Line 191: Couldn't send topology response received from", receivedFrom.GetAddress())
}

func GetP2PNetwork() []types.NetworkList {
	topologyChan := nodestate.InitTopologyChan()

	// Display the current P2P network connections
	topologyRequest := InitiateTopologyRequest()

	// Display waiting message
	// fmt.Println("\n\n\n...waiting")

	nodesFound := topologyRequest.NodesFound
	networkList := topologyRequest.ThisNodePeers

	listOfNetworkList := make([]types.NetworkList, 0)

	listOfNetworkList = append(listOfNetworkList, networkList)

	// log.Println("P2P: Before the loop")

	// log.Println("P2P: len(nodesFound.GetNodes()) =", len(nodesFound.GetNodes()))
	// log.Println("P2P: len(listOfNetworkList) =", len(listOfNetworkList))

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
			// Sleep for 500ms
			time.Sleep(500 * time.Millisecond)
		}
		// log.Println("P2P: len(nodesFound.GetNodes()) =", len(nodesFound.GetNodes()), "nodes: ", nodesFound.GetNodes())
		// log.Println("P2P: len(listOfNetworkList) =", len(listOfNetworkList))
		// for _, networkList := range listOfNetworkList {
		// 	// fmt.Println("P2P: networkList.Src =", networkList.Src.GetAddress())
		// }
	}

	// log.Println("P2P: After the loop")

	// close the channel
	nodestate.CloseTopologyChan()

	return listOfNetworkList
}
