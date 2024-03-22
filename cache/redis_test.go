package cache

//func TestRedisCache_Set(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	tests := []struct {
//		name       string
//		key        string
//		val        any
//		mock       func(ctrl *gomock.Controller)* redis.Cmdable
//		expiration time.Duration
//		wantErr    error
//	}{
//		// TODO: Add test cases.
//		{
//			name: "set value",
//			mock: func(ctrl *gomock.Controller) *redis.Cmdable {
//				client := mocks.NewMockCmdable(ctrl)
//				cmd := redis.NewCmd(context.Background())
//				cmd.SetVal("OK")
//
//				client.EXPECT().
//					Set(gomock.Any(), "key1", "value1", time.Minute).
//					Return(cmd)
//				return client
//			},
//			key:        "key1",
//			val:        "value1",
//			expiration: time.Second,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			c := NewRedisCache(tt.mock(ctrl))
//			err := c.Set(context.Background(), tt.key, tt.val, tt.expiration)
//			assert.Equal(t, tt.wantErr, err)
//		})
//	}
//}
