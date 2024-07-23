package net

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestConnect(t *testing.T) {
	go func() {
		err := Serve("tcp", ":8082")
		t.Log(err)
	}()
	time.Sleep(time.Second * 3)
	err := Connect("tcp", "localhost:8082")
	t.Log(err)
	//type args struct {
	//	network string
	//	address string
	//}
	//tests := []struct {
	//	name    string
	//	args    args
	//	wantErr assert.ErrorAssertionFunc
	//}{
	//	// TODO: Add test cases.
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		tt.wantErr(t, Connect(tt.args.network, tt.args.address), fmt.Sprintf("Connect(%v, %v)", tt.args.network, tt.args.address))
	//	})
	//}
}

func TestClient_Send(t *testing.T) {
	server := &Server{}
	go func() {
		err := server.Start("tcp", ":8081")
		t.Log(err)
	}()
	time.Sleep(time.Second * 3)
	client := &Client{
		network: "tcp",
		addr:    "localhost:8081",
	}
	resp, err := client.Send("hello")
	require.NoError(t, err)
	assert.Equal(t, "hellohello", resp)
}
