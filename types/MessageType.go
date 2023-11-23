package types

type MessageType uint8

const (
	MessageTypeTransaction MessageType = iota
	MessageTypeBlock
	MessageTypeConnectionRequest
	MessageTypeConnectionResponse
	MessageTypeBootstrapConnectionRequest
	MessageTypeBootstrapConnectionResponse
	MessageTypeBootstrapPingRequest
	MessageTypeBootstrapPingResponse
	MessageTypeUnknown
	MessageTypeFailure
)
