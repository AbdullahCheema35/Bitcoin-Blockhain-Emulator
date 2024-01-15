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
	MessageTypeBlockRequest
	MessageTypeBlockResponse
	MessageTypeBlockChainRequest
	MessageTypeBlockChainResponse
	MessageTypeTopologyRequest
	MessageTypeTopologyResponse
)
