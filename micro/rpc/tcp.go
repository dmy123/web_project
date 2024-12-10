package rpc

import (
	"encoding/binary"
	"net"
)

func ReadMsg(conn net.Conn) ([]byte, error) {
	lenBs := make([]byte, numOfLengthBytes)
	_, err := conn.Read(lenBs)
	if err != nil {
		return nil, err
	}
	headerLength := binary.BigEndian.Uint32(lenBs[:4])
	bodyLength := binary.BigEndian.Uint32(lenBs[4:])
	length := headerLength + bodyLength

	data := make([]byte, length)
	_, err = conn.Read(data[8:])
	copy(data[:8], lenBs)
	return data, err

}

//func EncodeMsg(data []byte) []byte {
//	reqLen := len(data)
//
//	// 我要在这，构建请求数据
//	// data = reqLen 的 64 位表示 + respData
//	res := make([]byte, reqLen+numOfLengthBytes)
//	// 第一步：
//	// 把长度写进去前八个字节
//	binary.BigEndian.PutUint64(res[:numOfLengthBytes], uint64(reqLen))
//	// 第二步：
//	// 写入数据
//	copy(res[numOfLengthBytes:], data)
//	return res
//}
