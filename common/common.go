// Some common functions

package common

import (
	"encoding/gob"
	"log"
	"net"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

func SendMessage(conn net.Conn, message types.Message) bool {
	enc := gob.NewEncoder(conn)
	err := enc.Encode(message)
	if err != nil {
		log.Println("Error encoding:", err)
		return false
	}
	return true
}

func ReceiveMessage(conn net.Conn) (bool, types.Message) {
	dec := gob.NewDecoder(conn)
	var message types.Message
	err := dec.Decode(&message)
	if err != nil {
		log.Println("Error decoding:", err)
		return false, types.Message{}
	}
	return true, message
}
