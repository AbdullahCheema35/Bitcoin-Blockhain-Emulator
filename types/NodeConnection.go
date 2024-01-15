package types

import (
	"encoding/gob"
	"net"
)

// NodeConnection represents a connection with a node
type NodeConnection struct {
	Node NodeAddress
	Conn net.Conn
	Enc  chan *gob.Encoder
	Dec  chan *gob.Decoder
}

// NewNodeConnection creates a new NodeConnection instance
func NewNodeConnection(node NodeAddress, conn net.Conn) NodeConnection {
	// Create channels
	enc_chan := make(chan *gob.Encoder)
	dec_chan := make(chan *gob.Decoder)

	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	// Initialize channels
	enc_chan <- encoder
	dec_chan <- decoder

	return NodeConnection{Node: node, Conn: conn, Enc: enc_chan, Dec: dec_chan}
}

func (nc *NodeConnection) SetNodeAddress(nodeAddr NodeAddress) {
	nc.Node = nodeAddr
}

func (nc *NodeConnection) GetEncoder() *gob.Encoder {
	return <-nc.Enc
}

func (nc *NodeConnection) GetDecoder() *gob.Decoder {
	return <-nc.Dec
}

func (nc *NodeConnection) SetEncoder(encoder *gob.Encoder) {
	nc.Enc <- encoder
}

func (nc *NodeConnection) SetDecoder(decoder *gob.Decoder) {
	nc.Dec <- decoder
}
