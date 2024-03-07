package additional

import (
	"reflect"
	"testing"
)

func TestNewMyStruct(t *testing.T) {
	type args struct {
		id   uint64
		name string
		opts []MyStructOption
	}
	tests := []struct {
		name string
		args args
		want MyStruct
	}{
		// TODO: Add test cases.
		{
			name: "normal",
			args: args{
				id:   1,
				name: "tom",
				opts: []MyStructOption{
					AddAddress("haha"),
				},
			},
			want: MyStruct{
				id:      1,
				name:    "tom",
				address: "haha",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMyStruct(tt.args.id, tt.args.name, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMyStruct() = %v, want %v", got, tt.want)
			}
		})
	}
}
