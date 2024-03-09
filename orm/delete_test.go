package orm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleter_Build(t *testing.T) {
	type testCase[T any] struct {
		name    string
		d       *Deleter[T]
		want    *Query
		wantErr error
	}
	tests := []testCase[TestModel]{
		// TODO: Add test cases.
		{
			name: "",
			d:    NewDeleter[TestModel](memoryDB(t)),
			want: &Query{
				SQL: "DELETE FROM `test_model`;",
			},
		},
		{
			name: "no where",
			d:    NewDeleter[TestModel](memoryDB(t)).From("`test_model`"),
			want: &Query{
				SQL: "DELETE FROM `test_model`;",
			},
		},
		{
			name: "where",
			d:    NewDeleter[TestModel](memoryDB(t)).From("`TestModel`").Where(C("Id").Eq(16)),
			want: &Query{
				SQL:  "DELETE FROM `TestModel` WHERE (`id`= ?);",
				Args: []any{16},
			},
		},
		{
			name: "from",
			d:    NewDeleter[TestModel](memoryDB(t)).From("`test_model`").Where(C("Id").Eq(16)),
			want: &Query{
				SQL:  "DELETE FROM `test_model` WHERE (`id`= ?);",
				Args: []any{16},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.Build()
			assert.Equal(t, err, tt.wantErr)
			if err != nil {
				return
			}
			assert.Equalf(t, tt.want, got, "Build()")
		})
	}
}
