package message

import (
	"bytes"
	"encoding/binary"
)

type Request struct {
	HeadLength  uint32
	BodyLength  uint32
	MessageId   uint32
	Version     uint8
	Compresser  uint8
	Serializer  uint8
	ServiceName string
	MethodName  string

	// 扩展字段，用于传递自定义元数据
	Meta map[string]string

	Data []byte
	//Data        any
	//Data        []any
}

func EncodeReq(request *Request) []byte {
	bs := make([]byte, request.HeadLength+request.BodyLength)

	binary.BigEndian.PutUint32(bs[:4], request.HeadLength)
	binary.BigEndian.PutUint32(bs[4:8], request.BodyLength)
	binary.BigEndian.PutUint32(bs[8:12], request.MessageId)
	bs[12] = request.Version
	bs[13] = request.Compresser
	bs[14] = request.Serializer
	cur := bs[15:]
	copy(cur, request.ServiceName)
	cur = cur[len(request.ServiceName):]
	cur[0] = '\n'
	cur = cur[1:]
	copy(cur, request.MethodName)
	cur = cur[len(request.MethodName):]
	cur[0] = '\n'
	cur = cur[1:]
	for key, value := range request.Meta {
		copy(cur, key)
		cur = cur[len(key):]
		cur[0] = '\r'
		cur = cur[1:]
		copy(cur, value)
		cur = cur[len(value):]
		cur[0] = '\n'
		cur = cur[1:]
	}

	copy(cur, request.Data)
	return bs
}

func DecodeReq(data []byte) *Request {
	resp := &Request{}
	resp.HeadLength = binary.BigEndian.Uint32(data[:4])
	resp.BodyLength = binary.BigEndian.Uint32(data[4:8])
	resp.MessageId = binary.BigEndian.Uint32(data[8:12])
	resp.Version = data[12]
	resp.Compresser = data[13]
	resp.Serializer = data[14]
	header := data[15:resp.HeadLength]
	// 用分隔符\n
	index := bytes.IndexByte(header, '\n')
	resp.ServiceName = string(header[:index])
	header = header[index+1:]
	index = bytes.IndexByte(header, '\n')
	resp.MethodName = string(header[:index])
	header = header[index+1:]
	index = bytes.IndexByte(header, '\n')
	if index != -1 {
		if resp.Meta == nil {
			resp.Meta = make(map[string]string)
		}
	}
	for index != -1 {
		pair := header[:index]
		pairIndex := bytes.IndexByte(pair, '\r')
		resp.Meta[string(pair[:pairIndex])] = string(pair[pairIndex+1:])
		header = header[index+1:]
		index = bytes.IndexByte(header, '\n')
	}
	if resp.BodyLength > 0 {
		resp.Data = data[resp.HeadLength:]
	}

	return resp
}

func (req *Request) CalculateHeadLength() {
	req.HeadLength = 15 + uint32(len(req.ServiceName)+1+len(req.MethodName)+1)
	for key, val := range req.Meta {
		req.HeadLength += uint32(len(key)+len(val)) + 2
	}
}

func (req *Request) CalculateBodyLength() {
	req.BodyLength = uint32(len(req.Data))
}
