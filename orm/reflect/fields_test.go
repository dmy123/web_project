package reflect

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestIterateFields(t *testing.T) {
	type User struct {
		Name string
		age  int
	}

	type args struct {
		entity any
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]any
		wantErr error
	}{
		{
			name: "struct",
			args: args{
				entity: User{
					Name: "tom",
					//age:  10,
				},
			},
			want: map[string]any{
				"Name": "tom",
				"age":  0, // 私有字段，拿不到，用0值填充
			},
		},
		{
			name: "pointer",
			args: args{
				entity: &User{
					Name: "tom",
					//age:  10,
				},
			},
			want: map[string]any{
				"Name": "tom",
				"age":  0, // 私有字段，拿不到，用0值填充
			},
		},
		{
			name: "basic type",
			args: args{
				entity: 18,
			},
			wantErr: errors.New("不支持类型"),
		},
		{
			name: "multiple pointer",
			args: args{
				entity: func() **User {
					res := &User{
						Name: "tom",
						//age:  10,
					}
					return &res
				}(),
			},
			want: map[string]any{
				"Name": "tom",
				"age":  0, // 私有字段，拿不到，用0值填充
			},
		},
		{
			name:    "nil",
			args:    args{},
			wantErr: errors.New("不支持 nil"),
		},
		{
			name: "User nil",
			args: args{
				entity: (*User)(nil),
			},
			wantErr: errors.New("不支持零值"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IterateFields(tt.args.entity)
			if err != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("IterateFields() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IterateFields() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetField(t *testing.T) {
	type User struct {
		Name string
		age  int
	}

	type args struct {
		entity   any
		field    string
		newValue any
	}
	tests := []struct {
		name       string
		args       args
		wantErr    error
		wantEntity any
	}{
		{
			name: "struct", // 结构体类型不可修改字段
			args: args{
				entity: User{
					Name: "Tom",
				},
				field:    "Name",
				newValue: "Jerry",
			},
			wantErr: errors.New("不可修改字段"),
		},
		{
			name: "pointer", // 指针类型可以修改
			args: args{
				entity: &User{
					Name: "Tom",
				},
				field:    "Name",
				newValue: "Jerry",
			},
			wantEntity: &User{
				Name: "Jerry",
			},
		},
		{
			name: "pointer exported", // 未导出字段不可修改
			args: args{
				entity: &User{
					Name: "Tom",
					age:  10,
				},
				field:    "age",
				newValue: "3",
			},
			wantEntity: &User{
				Name: "Tom",
				age:  10,
			},
			wantErr: errors.New("不可修改字段"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetField(tt.args.entity, tt.args.field, tt.args.newValue); err != nil {
				if tt.wantErr == nil || err.Error() != tt.wantErr.Error() {
					t.Errorf("SetField() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if tt.wantEntity != nil {
				assert.Equal(t, tt.args.entity, tt.wantEntity)
			}
		})
	}

	var i = 0
	//reflect.ValueOf(i).Set(reflect.ValueOf(12))
	//assert.Equal(t, i, 12)
	ptr := &i
	reflect.ValueOf(ptr).Elem().Set(reflect.ValueOf(12))
	assert.Equal(t, i, 12)
}
