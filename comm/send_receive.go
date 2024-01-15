package comm

import (
	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

func SendMessage(nc types.NodeConnection, message types.Message) bool {
	enc := nc.GetEncoder()
	defer func() {
		nc.SetEncoder(enc)
	}()
	err := enc.Encode(&message)
	if err != nil {
		// log.Println("Error encoding:", err)
		// log.Panicf("Error encoding message of type (%d): %v\n", message.Header.Type, err)
		return false
	}
	return true
}

func ReceiveMessage(nc types.NodeConnection) (bool, types.Message) {
	dec := nc.GetDecoder()
	defer func() {
		nc.SetDecoder(dec)
	}()
	var message types.Message
	err := dec.Decode(&message)
	if err != nil {
		// log.Println("Error decoding:", err)
		// log.Panicln("Error decoding:", err)
		return false, message
	}
	return true, message
}
