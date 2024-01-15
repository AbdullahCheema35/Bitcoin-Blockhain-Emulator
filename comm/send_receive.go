package comm

import (
	"encoding/gob"
	"log"
	"net"

	"github.com/AbdullahCheema35/Bitcoin-Blockhain-Emulator/types"
)

func SendMessage(conn net.Conn, message types.Message) bool {
	enc := gob.NewEncoder(conn)
	err := enc.Encode(&message)
	if err != nil {
		// log.Println("Error encoding:", err)
		log.Panicf("Error encoding message of type (%d): %v\n", message.Header.Type, err)
		return false
	}
	return true
}

func ReceiveMessage(conn net.Conn) (bool, types.Message) {
	dec := gob.NewDecoder(conn)
	var message types.Message
	err := dec.Decode(&message)
	if err != nil {
		// log.Println("Error decoding:", err)
		log.Panicln("Error decoding:", err)
		return false, message
	}
	return true, message
}
