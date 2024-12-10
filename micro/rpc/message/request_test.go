package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeDecodeReq(t *testing.T) {
	tests := []struct {
		name    string
		request *Request
	}{
		{
			name: "TestEncodeDecodeReq",
			request: &Request{
				MessageId:   123,
				Version:     12,
				Compresser:  13,
				Serializer:  14,
				ServiceName: "user-service",
				MethodName:  "GetById",
				Meta: map[string]string{
					"trace-id": "123456",
					"a/b":      "a",
				},
				Data: []byte("hello world"),
			},
		},
		{
			name: "TestEncodeDecodeReq, without meta",
			request: &Request{
				MessageId:   123,
				Version:     12,
				Compresser:  13,
				Serializer:  14,
				ServiceName: "user-service",
				MethodName:  "GetById",
				Data:        []byte("hello world"),
			},
		},
		{
			name: "TestEncodeDecodeReq, without meta and data",
			request: &Request{
				MessageId:   123,
				Version:     12,
				Compresser:  13,
				Serializer:  14,
				ServiceName: "user-service",
				MethodName:  "GetById",
			},
		},
		{
			name: "TestEncodeDecodeReq, data with \n",
			request: &Request{
				MessageId:   123,
				Version:     12,
				Compresser:  13,
				Serializer:  14,
				ServiceName: "user-service",
				MethodName:  "GetById",
				Data:        []byte("hello \n world"),
			},
		},
		// 禁止开发者在header使用\n\r
		//{
		//	name: "TestEncodeDecodeReq, ServiceName with \n",
		//	request: &Request{
		//		MessageId:   123,
		//		Version:     12,
		//		Compresser:  13,
		//		Serializer:  14,
		//		ServiceName: "user\n\r-service",
		//		MethodName:  "GetById",
		//		Data:        []byte("hello \n world"),
		//	},
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.request.CalculateHeadLength()
			tt.request.CalculateBodyLength()
			data := EncodeReq(tt.request)
			req := DecodeReq(data)
			assert.Equal(t, tt.request, req)
		})
	}
}
