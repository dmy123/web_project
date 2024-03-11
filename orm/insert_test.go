package orm

import (
	"awesomeProject1/orm/internal/errs"
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInserter_SQLite_upsert(t *testing.T) {
	db := memoryDBOpt(t, DBWithDialect(DialectSQLite))
	type testCase[T any] struct {
		name    string
		i       *Inserter[TestModel]
		want    *Query
		wantErr error
	}
	tests := []testCase[TestModel]{
		//{
		//	name: "",
		//	i: NewInserter[TestModel](db).Values().Upsert().ConflictColumns().Update(),
		//	want: ,
		//},
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
			}).OnDuplicateKey().ConflictColumns("Id").Update(Assign("FirstName", "Deng"), Assign("Age", 19)),
			want: &Query{SQL: "INSERT INTO `test_model` (`id`, `first_name`, `age`, `last_name`) VALUES (?,?,?,?) ON CONFLICT(`id`) DO UPDATE SET `first_name`=?, `age`=?;",
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
			}).OnDuplicateKey().ConflictColumns("FirstName", "LastName").Update(C("FirstName"), C("Age")),
			want: &Query{SQL: "INSERT INTO `test_model` (`id`, `first_name`, `age`, `last_name`) VALUES (?,?,?,?) ON CONFLICT(`first_name`, `last_name`) DO UPDATE SET `first_name`=excluded.`first_name`, `age`=excluded.`age`;",
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

func TestInserter_Exec(t *testing.T) {
	//var i Inserter[TestModel]
	//// 判断两次
	////res, err := i.Exec(context.Background())
	////if err != nil {
	////
	////}
	////affectedRows, err := res.RowsAffected()
	////if err != nil {
	////
	////}
	//// 封装后
	//res := i.Exec(context.Background())
	//result, err := res.RowsAffected()
	//if err != nil {
	//	return
	//}
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := OpenDB(mockDB)
	require.NoError(t, err)

	testCases := []struct {
		name     string
		i        *Inserter[TestModel]
		wantErr  error
		affected int64
	}{
		{
			name: "query error",
			i: func() *Inserter[TestModel] {
				return NewInserter[TestModel](db).Values(&TestModel{}).Columns("invalid")
			}(),
			wantErr: errs.NewErrUnknownField("invalid"),
		},
		{
			name: "db error",
			i: func() *Inserter[TestModel] {
				mock.ExpectExec("INSERT INTO .*").WillReturnError(errors.New("db error"))
				return NewInserter[TestModel](db).Values(&TestModel{})
			}(),
			wantErr: errors.New("db error"),
		},
		{
			name: "db",
			i: func() *Inserter[TestModel] {
				res := driver.RowsAffected(1)
				mock.ExpectExec("INSERT INTO .*").WillReturnResult(res)
				return NewInserter[TestModel](db).Values(&TestModel{})
			}(),
			affected: 1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.i.Exec(context.Background())
			af, err := res.RowsAffected()
			assert.Equal(t, tc.wantErr, res.err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.affected, af)
		})
	}
}
