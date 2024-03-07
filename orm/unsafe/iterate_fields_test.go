package unsafe

import "testing"

func TestPrintFieldOffset(t *testing.T) {
	type args struct {
		entity any
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "User",
			args: args{
				entity: User{},
			},
		},
		{
			name: "UserV1",
			args: args{
				entity: UserV1{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintFieldOffset(tt.args.entity)
		})
	}
}

type User struct {
	//0
	Name string
	//16
	Age int32
	//24
	Alias []string
	//48
	Address string
}

type UserV1 struct {
	//0
	Name string
	//16
	Age int32
	// 20
	Agev1 int32
	//24
	Alias []string
	//48
	Address string
}
