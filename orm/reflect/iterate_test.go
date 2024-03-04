package reflect

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIterateArray(t *testing.T) {
	type args struct {
		entity any
	}
	tests := []struct {
		name    string
		args    args
		want    []any
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "array",
			args: args{
				entity: [3]int{1, 2, 3},
			},
			want: []any{
				1, 2, 3,
			},
		},
		{
			name: "slice",
			args: args{
				entity: []int{1, 2, 3},
			},
			want: []any{
				1, 2, 3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IterateArrayOrSlice(tt.args.entity)
			if err != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("IterateArrayOrSlice() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			assert.Equalf(t, tt.want, got, "IterateArrayOrSlice(%v)", tt.args.entity)
		})
	}
}

func TestIterateMap(t *testing.T) {
	type args struct {
		entity any
	}
	tests := []struct {
		name    string
		args    args
		want    []any
		want1   []any
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "map",
			args: args{
				entity: map[string]string{
					"A": "a",
					"B": "b",
				},
			},
			want: []any{
				"A", "B",
			},
			want1: []any{
				"a", "b",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := IterateMap(tt.args.entity)
			if err != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("IterateMap() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			assert.EqualValuesf(t, tt.want, got, "IterateMap(%v)", tt.args.entity)
			assert.EqualValuesf(t, tt.want1, got1, "IterateMap(%v)", tt.args.entity)
			//assert.Equalf(t, tt.want, got, "IterateMap(%v)", tt.args.entity)
			//assert.Equalf(t, tt.want1, got1, "IterateMap(%v)", tt.args.entity)
		})
	}
}
