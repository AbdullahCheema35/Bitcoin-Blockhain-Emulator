package types

type Message struct {
	Header MessageHeader
	Body   interface{}
}

func NewMessage(header MessageHeader, body interface{}) Message {
	return Message{Header: header, Body: body}
}
