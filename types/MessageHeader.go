package types

type MessageHeader struct {
	Type   MessageType
	Sender NodeAddress
}

func NewMessageHeader(messageType MessageType, sender NodeAddress) MessageHeader {
	return MessageHeader{Type: messageType, Sender: sender}
}
