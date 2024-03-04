package reflect

import (
	"awesomeProject1/orm/reflect/types"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestIterateFunc(t *testing.T) {
	type args struct {
		entity any
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]FuncInfo
		wantErr error
	}{
		// TODO: Add test cases.
		//以结构体作为输入，只能访问到结构体作为接收器的方法
		//以指针作为输入，能访问到任何接收器方法
		{
			name: "struct",
			args: args{
				entity: types.NewUser("Tom", 18),
			},
			want: map[string]FuncInfo{
				"GetAge": {
					Name:       "GetAge",
					InputTypes: []reflect.Type{reflect.TypeOf(types.User{})},
					OutputTypes: []reflect.Type{
						reflect.TypeOf(0),
					},
					Result: []any{18},
				},
				//"ChangeName": {
				//	Name: "ChangeName",
				//	InputTypes: []reflect.Type{
				//		reflect.TypeOf(""),
				//	},
				//},
			},
		},
		{
			name: "pointer",
			args: args{
				entity: types.NewUserPtr("Tom", 18),
			},
			want: map[string]FuncInfo{
				"GetAge": {
					Name:       "GetAge",
					InputTypes: []reflect.Type{reflect.TypeOf(&types.User{})},
					OutputTypes: []reflect.Type{
						reflect.TypeOf(0),
					},
					Result: []any{18},
				},
				"ChangeName": {
					Name:        "ChangeName",
					InputTypes:  []reflect.Type{reflect.TypeOf(&types.User{}), reflect.TypeOf("")},
					OutputTypes: []reflect.Type{},
					Result:      []any{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IterateFunc(tt.args.entity)
			if err != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("IterateFunc() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			assert.Equalf(t, tt.want, got, "IterateFunc(%v)", tt.args.entity)
		})
	}
}
