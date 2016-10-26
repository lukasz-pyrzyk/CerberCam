package main

import proto "github.com/golang/protobuf/proto"

// Deserialize
func Deserialize(body []byte) *Message {
	log.Debugf("Deserializing message with %d bytes", len(body))

	msg := &Message{}
	err := proto.Unmarshal(body, msg)

	if err != nil {
		log.Criticalf("%s: %s", msg, err)
	}
	return msg
}
