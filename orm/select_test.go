package orm

import (
	"awesomeProject1/orm/internal/errs"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelector_Build(t *testing.T) {
	//db, err := NewDB()
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
			builder: NewSelector[TestModel](MustNewDB()),
			wantQuery: &Query{
				//SQL:  "SELECT * FROM `TestModel`;",
				SQL:  "SELECT * FROM `test_model`;",
				Args: nil,
			},
		},
		{
			name:    "from",
			builder: NewSelector[TestModel](MustNewDB()).From("`TestModel`"),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel`;",
				Args: nil,
			},
		},
		{
			name: "empty from",
			//builder: (&Selector[TestModel]{}).From(""),
			builder: NewSelector[TestModel](MustNewDB()).From(""),
			wantQuery: &Query{
				//SQL:  "SELECT * FROM `TestModel`;",
				Args: nil,
				SQL:  "SELECT * FROM `test_model`;",
			},
		},
		{
			name:    "from db",
			builder: NewSelector[TestModel](MustNewDB()).From("`test_db`.`test_model`"),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_db`.`test_model`;",
				Args: nil,
			},
		},
		{
			name:    "empty where",
			builder: NewSelector[TestModel](MustNewDB()).Where(),
			wantQuery: &Query{
				//SQL: "SELECT * FROM `TestModel`;",
				SQL: "SELECT * FROM `test_model`;",
				//Args: []any{123},
			},
		},
		{
			name:    "where",
			builder: NewSelector[TestModel](MustNewDB()).Where(C("Age").Eq(123)),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` WHERE (`age` = ?);",
				//SQL:  "SELECT * FROM `TestModel` WHERE (`Age` = ?);",
				Args: []any{123},
			},
		},
		{
			name:    "not",
			builder: NewSelector[TestModel](MustNewDB()).Where(Not(C("Age").Eq(123))),
			wantQuery: &Query{
				//SQL:  "SELECT * FROM `TestModel` WHERE (NOT (`Age` = ?));",
				SQL:  "SELECT * FROM `test_model` WHERE (NOT (`age` = ?));",
				Args: []any{123},
			},
		},
		{
			name:    "not",
			builder: NewSelector[TestModel](MustNewDB()).Where(Not(C("Age").Eq(123)).And(C("Id").Eq(321))),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model` WHERE ((NOT (`age` = ?))AND (`id` = ?));",
				//SQL:  "SELECT * FROM `TestModel` WHERE ((NOT (`Age` = ?))AND (`Id` = ?));",
				Args: []any{123, 321},
			},
		},
		{
			name:    "invalid column",
			builder: NewSelector[TestModel](MustNewDB()).Where(Not(C("haha").Eq(123)).And(C("Id").Eq(321))),
			wantErr: errs.NewErrUnknownField("haha"),
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
