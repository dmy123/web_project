package proto

import (
	"errors"
	"google.golang.org/protobuf/proto"
)

type Serializer struct {
}

func (s Serializer) Code() byte {
	return 2
}

func (s Serializer) Encode(val any) ([]byte, error) {
	msg, ok := val.(proto.Message)
	if !ok {
		return nil, errors.New("proto.Message is not protobuf")
	}
	return proto.Marshal(msg)
}

func (s Serializer) Decode(data []byte, val any) error {
	msg, ok := val.(proto.Message)
	if !ok {
		return errors.New("proto.Message is not protobuf")
	}
	return proto.Unmarshal(data, msg)
}
