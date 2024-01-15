package comm

import (
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

func CreateMessage(sender types.NodeAddress, msgType types.MessageType, body interface{}) types.Message {
	messageHeader := types.NewMessageHeader(msgType, sender)
	message := types.NewMessage(messageHeader, body)
	return message
}
