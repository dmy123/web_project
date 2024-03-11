package orm

import (
	"awesomeProject1/orm/internal/errs"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInserter_Build(t *testing.T) {
	db := memoryDB(t)
	type testCase[T any] struct {
		name    string
		i       *Inserter[TestModel]
		want    *Query
		wantErr error
	}
	tests := []testCase[TestModel]{
		// TODO: Add test cases.
		{
			name: "single row",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id:        12,
				FirstName: "Tom",
				Age:       8,
				LastName: &sql.NullString{
					String: "Jerry",
					Valid:  true,
				},
			}),
			want: &Query{SQL: "INSERT INTO `test_model` (`id`, `first_name`, `age`, `last_name`) VALUES (?,?,?,?);",
				Args: []any{int64(12), "Tom", int8(8),
					&sql.NullString{String: "Jerry", Valid: true}}},
		},
		{
			name: "multiple row",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id:        12,
				FirstName: "Tom",
				Age:       8,
				LastName: &sql.NullString{
					String: "Jerry",
					Valid:  true,
				},
			}, &TestModel{
				Id:        2,
				FirstName: "sam",
				Age:       18,
				LastName: &sql.NullString{
					String: "harry",
					Valid:  true,
				},
			}),
			want: &Query{SQL: "INSERT INTO `test_model` (`id`, `first_name`, `age`, `last_name`) VALUES (?,?,?,?),(?,?,?,?);",
				Args: []any{int64(12), "Tom", int8(8), &sql.NullString{String: "Jerry", Valid: true},
					int64(2), "sam", int8(18), &sql.NullString{String: "harry", Valid: true}}},
		},
		{
			name:    "no row",
			i:       NewInserter[TestModel](db).Values(),
			wantErr: errs.ErrInsertZeroRows,
		},
		{
			name: "partial columns",
			i: NewInserter[TestModel](db).Columns("Id", "FirstName").Values(&TestModel{
				Id:        12,
				FirstName: "Tom",
				Age:       8,
				LastName: &sql.NullString{
					String: "Jerry",
					Valid:  true,
				},
			}),
			want: &Query{SQL: "INSERT INTO `test_model` (`id`, `first_name`) VALUES (?,?);",
				Args: []any{int64(12), "Tom"}},
		},
		{
			name: "multiple row",
			i: NewInserter[TestModel](db).Columns("Id", "FirstName").Values(&TestModel{
				Id:        12,
				FirstName: "Tom",
				Age:       8,
				LastName: &sql.NullString{
					String: "Jerry",
					Valid:  true,
				},
			}, &TestModel{
				Id:        2,
				FirstName: "sam",
				Age:       18,
				LastName: &sql.NullString{
					String: "harry",
					Valid:  true,
				},
			}),
			want: &Query{SQL: "INSERT INTO `test_model` (`id`, `first_name`) VALUES (?,?),(?,?);",
				Args: []any{int64(12), "Tom", int64(2), "sam"}},
		},
		{
			name: "upsert",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id:        12,
				FirstName: "Tom",
				Age:       8,
				LastName: &sql.NullString{
					String: "Jerry",
					Valid:  true,
				},
			}).OnDuplicateKey().Update(Assign("FirstName", "Deng"), Assign("Age", 19)),
			want: &Query{SQL: "INSERT INTO `test_model` (`id`, `first_name`, `age`, `last_name`) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE `first_name`=?, `age`=?;",
				Args: []any{int64(12), "Tom", int8(8),
					&sql.NullString{String: "Jerry", Valid: true}, "Deng", 19}},
		},
		{
			name: "upsert column",
			i: NewInserter[TestModel](db).Values(&TestModel{
				Id:        12,
				FirstName: "Tom",
				Age:       8,
				LastName: &sql.NullString{
					String: "Jerry",
					Valid:  true,
				},
			}).OnDuplicateKey().Update(C("FirstName"), C("Age")),
			want: &Query{SQL: "INSERT INTO `test_model` (`id`, `first_name`, `age`, `last_name`) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE `first_name`=VALUES(`first_name`), `age`=VALUES(`age`);",
				Args: []any{int64(12), "Tom", int8(8),
					&sql.NullString{String: "Jerry", Valid: true}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.Build()
			assert.Equal(t, err, tt.wantErr)
			if err != nil {
				return
			}
			assert.Equalf(t, tt.want, got, "Build()")
		})
	}
}
