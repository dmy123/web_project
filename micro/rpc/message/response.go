package message

import (
	"encoding/binary"
)

type Response struct {
	HeadLength uint32
	BodyLength uint32
	MessageId  uint32
	Version    uint8
	Compresser uint8
	Serializer uint8
	Error      []byte
	Data       []byte
}

func EncodeResp(response *Response) []byte {
	bs := make([]byte, response.HeadLength+response.BodyLength)

	binary.BigEndian.PutUint32(bs[:4], response.HeadLength)
	binary.BigEndian.PutUint32(bs[4:8], response.BodyLength)
	binary.BigEndian.PutUint32(bs[8:12], response.MessageId)
	bs[12] = response.Version
	bs[13] = response.Compresser
	bs[14] = response.Serializer
	copy(bs[15:response.HeadLength], response.Error)
	copy(bs[response.HeadLength:], response.Data)
	return bs
}

func DecodeResp(data []byte) *Response {
	resp := &Response{}
	if len(data) == 0 {
		return resp
	}
	resp.HeadLength = binary.BigEndian.Uint32(data[:4])
	resp.BodyLength = binary.BigEndian.Uint32(data[4:8])
	resp.MessageId = binary.BigEndian.Uint32(data[8:12])
	resp.Version = data[12]
	resp.Compresser = data[13]
	resp.Serializer = data[14]
	if resp.HeadLength > 15 {
		resp.Error = data[15:resp.HeadLength]
	}
	if resp.BodyLength > 0 {
		resp.Data = data[resp.HeadLength:]
	}

	return resp
}

func (r *Response) CalculateHeadLength() {
	r.HeadLength = 15 + uint32(len(r.Error))
}

func (r *Response) CalculateBodyLength() {
	r.BodyLength = uint32(len(r.Data))
}
