package message

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncodeDecodeResp(t *testing.T) {
	tests := []struct {
		name     string
		response *Response
	}{
		{
			name: "TestEncodeDecodeResp",
			response: &Response{
				MessageId:  123,
				Version:    12,
				Compresser: 13,
				Serializer: 14,
				Error:      []byte("test error"),
				Data:       []byte("hello world"),
			},
		},
		{
			name: "TestEncodeDecodeResp, without error",
			response: &Response{
				MessageId:  123,
				Version:    12,
				Compresser: 13,
				Serializer: 14,
				Data:       []byte("hello world"),
			},
		},
		{
			name: "TestEncodeDecodeResp, without data",
			response: &Response{
				MessageId:  123,
				Version:    12,
				Compresser: 13,
				Serializer: 14,
				Error:      []byte("test error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.response.CalculateHeadLength()
			tt.response.CalculateBodyLength()
			data := EncodeResp(tt.response)
			req := DecodeResp(data)
			assert.Equal(t, tt.response, req)
		})
	}
}
