package orm

import (
	"awesomeProject1/orm/internal/errs"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSelector_Build(t *testing.T) {
	//db, err := Open()
	//require.NoError(t, err)
	//type testCase[T any] struct {
	//	name    string
	//	s       Selector[T]
	//	want    *Query
	//	wantErr bool
	//}
	//tests := []testCase[ /* TODO: Insert concrete types here */ ]{
	//	// TODO: Add test cases.
	//}
	tests := []struct {
		name string
		//wantErr bool
		//want bool
		builder   QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name: "no from",
			//builder: &Selector[TestModel]{},
			builder: NewSelector[TestModel](memoryDB(t)),
			wantQuery: &Query{
				//SQL:  "SELECT * FROM `TestModel`;",
				SQL:  "SELECT * FROM `test_model`;",
				Args: nil,
			},
		},
		{
			name:    "from",
			builder: NewSelector[TestModel](memoryDB(t)).From("`TestModel`"),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel`;",
				Args: nil,
			},
		},
		{
			name: "empty from",
			//builder: (&Selector[TestModel]{}).From(""),
			builder: NewSelector[TestModel](memoryDB(t)).From(""),
			wantQuery: &Query{
				//SQL:  "SELECT * FROM `TestModel`;",
				Args: nil,
				SQL:  "SELECT * FROM `test_model`;",
			},
		},
		{
			name:    "from db",
			builder: NewSelector[TestModel](memoryDB(t)).From("`test_db`.`test_model`"),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_db`.`test_model`;",
				Args: nil,
			},
		},
		{
			name:    "empty where",
			builder: NewSelector[TestModel](memoryDB(t)).Where(),
			wantQuery: &Query{
				//SQL: "SELECT * FROM `TestModel`;",
				SQL: "SELECT * FROM `test_model`;",
				//Args: []any{123},
			},
		},
		{
			name:    "where",
			builder: NewSelector[TestModel](memoryDB(t)).Where(C("Age").Eq(123)),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` WHERE (`age`= ?);",
				//SQL:  "SELECT * FROM `TestModel` WHERE (`Age` = ?);",
				Args: []any{123},
			},
		},
		{
			name:    "not",
			builder: NewSelector[TestModel](memoryDB(t)).Where(Not(C("Age").Eq(123))),
			wantQuery: &Query{
				//SQL:  "SELECT * FROM `TestModel` WHERE (NOT (`Age` = ?));",
				SQL:  "SELECT * FROM `test_model` WHERE (NOT (`age`= ?));",
				Args: []any{123},
			},
		},
		{
			name:    "not",
			builder: NewSelector[TestModel](memoryDB(t)).Where(Not(C("Age").Eq(123)).And(C("Id").Eq(321))),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` WHERE ((NOT (`age`= ?))AND (`id`= ?));",
				//SQL:  "SELECT * FROM `TestModel` WHERE ((NOT (`Age` = ?))AND (`Id` = ?));",
				Args: []any{123, 321},
			},
		},
		{
			name:    "invalid column",
			builder: NewSelector[TestModel](memoryDB(t)).Where(Not(C("haha").Eq(123)).And(C("Id").Eq(321))),
			wantErr: errs.NewErrUnknownField("haha"),
		},
		{
			name:    "raw expression as predicate",
			builder: NewSelector[TestModel](memoryDB(t)).Where(Raw("`id`<?", 12).AsPredicate()),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE ((`id`<?));",
				Args: []any{12},
			},
		},
		{
			name:    "raw expression used in predicate",
			builder: NewSelector[TestModel](memoryDB(t)).Where(C("Id").Eq(Raw("`age`+ ?", 1))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`id`= (`age`+ ?));",
				Args: []any{1},
			},
		},
		{
			name:    "columns alias",
			builder: NewSelector[TestModel](memoryDB(t)).Where(C("Id").As("my_id").Eq(18)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`id`= ?);",
				Args: []any{18},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := tt.builder.Build()
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantQuery, q)
			//got, err := tt.s.Build()
			//if (err != nil) != tt.wantErr {
			//	t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
			//	return
			//}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("Build() got = %v, want %v", got, tt.want)
			//}
		})
	}
}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}

func memoryDB(T *testing.T) *DB {
	db, err := Open("sqlite", "file:test.db?cache=shared&mode=memory")
	require.NoError(T, err)
	return db
}

func TestSelector_Get(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	db, err := OpenDB(mockDB)
	require.NoError(t, err)

	mock.ExpectQuery("SELECT .*").WillReturnError(errors.New("query error"))

	rows := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	//data
	rows = sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
	rows.AddRow("1", "Tom", "18", "Jerry")
	mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	//scan err
	//rows = sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
	//rows.AddRow("ab", "Tom", "18", "Jerry")
	//mock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	fmt.Println(mock)
	type args struct {
		ctx context.Context
	}
	type testCase[T any] struct {
		name string
		s    *Selector[T]
		//args    args
		want    *T
		wantErr error
	}
	tests := []testCase[TestModel]{
		// TODO: Add test cases.
		//{
		//	name:    "invalid query",
		//	s:       NewSelector[TestModel](db).Where(C("XXX").Eq(1)),
		//	wantErr: errs.NewErrUnknownField("XXX"),
		//},
		{
			name:    "query error",
			s:       NewSelector[TestModel](db).Where(C("Id").Eq(1)),
			wantErr: errors.New("query error"),
		},
		{
			name:    "no rows",
			s:       NewSelector[TestModel](db).Where(C("Id").Lt(1)),
			wantErr: ErrNoRows,
		},
		{
			name: "data",
			s:    NewSelector[TestModel](db).Where(C("Id").Lt(1)),
			want: &TestModel{
				Id:        1,
				Age:       18,
				FirstName: "Tom",
				LastName: &sql.NullString{
					String: "Jerry",
					Valid:  true,
				},
			},
		},
		//{
		//	name:    "scan err",
		//	s:       NewSelector[TestModel](db).Where(C("Id").Lt(1)),
		//	wantErr: errs.ErrNoRows,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Get(context.Background())
			if err != nil || tt.wantErr != nil {
				assert.Equal(t, err, tt.wantErr)
				return
			}
			assert.Equalf(t, tt.want, got, "Get(%v)", context.Background())
		})
	}
}

func TestSelector_Select(t *testing.T) {
	db := memoryDB(t)
	type args struct {
		cols []string
	}
	type testCase struct {
		name    string
		s       QueryBuilder
		args    args
		want    *Query
		wantErr error
	}
	tests := []testCase{
		//{
		//	name: "multiple columns",
		//	s:    NewSelector[TestModel](db).Select("first_name", "last_name"),
		//	want: &Query{
		//		SQL: "SELECT `first_name`, `last_name` FROM `test_model`;",
		//	},
		//},
		{
			name: "multiple columns",
			s:    NewSelector[TestModel](db).Select(C("FirstName"), C("LastName")),
			want: &Query{
				SQL: "SELECT `first_name`, `last_name` FROM `test_model`;",
			},
		},
		{
			name:    "invalid columns",
			s:       NewSelector[TestModel](db).Select(C("Invalid")),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
		{
			name: "avg",
			s:    NewSelector[TestModel](db).Select(Avg("Age")),
			want: &Query{
				SQL: "SELECT AVG(`age`) FROM `test_model`;",
			},
		},
		{
			name: "sum",
			s:    NewSelector[TestModel](db).Select(Sum("Age")),
			want: &Query{
				SQL: "SELECT SUM(`age`) FROM `test_model`;",
			},
		},
		{
			name: "count",
			s:    NewSelector[TestModel](db).Select(Count("Age")),
			want: &Query{
				SQL: "SELECT COUNT(`age`) FROM `test_model`;",
			},
		},
		{
			name: "max",
			s:    NewSelector[TestModel](db).Select(Max("Age")),
			want: &Query{
				SQL: "SELECT MAX(`age`) FROM `test_model`;",
			},
		},
		{
			name: "min",
			s:    NewSelector[TestModel](db).Select(Min("Age")),
			want: &Query{
				SQL: "SELECT MIN(`age`) FROM `test_model`;",
			},
		},
		{
			name:    "aggregate invalid columns",
			s:       NewSelector[TestModel](db).Select(Min("Invalid")),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
		{
			name:    "multiple aggregate",
			s:       NewSelector[TestModel](db).Select(Min("Invalid"), Max("Age")),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
		{
			name: "raw expression",
			s:    NewSelector[TestModel](db).Select(Raw("COUNT(DISTINCT `first_name`)")),
			want: &Query{
				SQL: "SELECT COUNT(DISTINCT `first_name`) FROM `test_model`;",
			},
		},
		{
			name: "columns alias",
			s:    NewSelector[TestModel](db).Select(C("FirstName").As("my_name"), C("LastName")),
			want: &Query{
				SQL: "SELECT `first_name` AS `my_name`, `last_name` FROM `test_model`;",
			},
		},
		{
			name: "avg",
			s:    NewSelector[TestModel](db).Select(Avg("Age").As("age")),
			want: &Query{
				SQL: "SELECT AVG(`age`) AS `age` FROM `test_model`;",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := tt.s.Build()
			assert.Equal(t, err, tt.wantErr)
			if err != nil {
				return
			}
			assert.Equal(t, tt.want, q)
		})
	}
}
