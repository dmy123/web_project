package cache_api

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"
)

func TestBuildInMapCache_Get(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		cache   func() *BuildInMapCache
		want    any
		wantErr error
	}{
		{
			name: "key not found",
			key:  "not exist key",
			cache: func() *BuildInMapCache {
				return NewBuildInMapCache(time.Second * 10)
			},
			wantErr: fmt.Errorf("%w, key: %s", errKeyNotFound, "not exist key"),
		},
		{
			name: "get value",
			key:  "key1",
			cache: func() *BuildInMapCache {
				res := NewBuildInMapCache(time.Second * 10)
				err := res.Set(context.Background(), "key1", "val1", time.Second*10)
				require.NoError(t, err)
				return res
			},
			want: "val1",
		},
		{
			name: "expired",
			key:  "expired key",
			cache: func() *BuildInMapCache {
				res := NewBuildInMapCache(time.Second * 5)
				err := res.Set(context.Background(), "expired key", "val1", time.Second*1)
				require.NoError(t, err)
				time.Sleep(time.Second * 2)
				return res
			},
			wantErr: fmt.Errorf("%w, key: %s", errKeyNotFound, "expired key"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cache().Get(context.Background(), tt.key)
			assert.Equal(t, err, tt.wantErr)
			if err != nil {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildInMapCache_Loop(t *testing.T) {
	cnt := 0
	c := NewBuildInMapCache(time.Second, WithEvictedCallback(func(key string, val any) {
		cnt++
	}))
	err := c.Set(context.Background(), "key1", 111, time.Second)
	require.NoError(t, err)
	time.Sleep(time.Second * 3)
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	_, ok := c.data["key1"]
	require.False(t, ok)
	require.Equal(t, 1, cnt)
}
