package main

import proto "github.com/golang/protobuf/proto"

type serializer struct {
}

// Deserialize
func (s serializer) Deserialize(body []byte) *Message {
	log.Debugf("Deserializing message with %d bytes", len(body))

	msg := &Message{}
	err := proto.Unmarshal(body, msg)

	failOnError(err, "Cannot deserialize protobuf object")

	return msg
}

func (s serializer) Serialize(response *Response) (body []byte) {
	log.Debugf("Deserializing message with %d bytes", len(body))
	data, err := proto.Marshal(response)

	failOnError(err, "Cannot serialize protobuf object")

	return data
}
