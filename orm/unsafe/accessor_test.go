package unsafe

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestUnsafeAccessor_Field(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}
	//type fields struct {
	//	fields  map[string]FieldMeta
	//	address unsafe.Pointer
	//}
	type args struct {
		entity any
		field  string
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				entity: &User{
					Name: "tom",
					Age:  4,
				},
				field: "Age",
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewUnsafeAccessor(tt.args.entity)
			got, err := a.Field(tt.args.field)
			if (err != nil) != tt.wantErr {
				t.Errorf("Field() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Field() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnsafeAccessor_SetField(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}
	type args struct {
		entity any
		field  string
		val    any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    any
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				entity: &User{
					Name: "tom",
					Age:  4,
				},
				field: "Age",
				val:   18,
			},
			want: 18,
		},
		{
			name: "",
			args: args{
				entity: &User{
					Name: "tom",
					Age:  4,
				},
				field: "Name",
				val:   "jerry",
			},
			want: "jerry",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewUnsafeAccessor(tt.args.entity)
			if err := a.SetField(tt.args.field, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("SetField() error = %v, wantErr %v", err, tt.wantErr)
			}
			data, err := a.Field(tt.args.field)
			assert.NoError(t, err)
			assert.Equal(t, data, tt.want)
		})
	}
}
