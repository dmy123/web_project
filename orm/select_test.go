package orm

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelector_Build(t *testing.T) {
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
			name:    "no from",
			builder: &Selector[TestModel]{},
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel`;",
				Args: nil,
			},
		},
		{
			name:    "from",
			builder: &Selector[TestModel]{table: "`TestModel`"},
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel`;",
				Args: nil,
			},
		},
		{
			name:    "empty from",
			builder: (&Selector[TestModel]{}).From(""),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel`;",
				Args: nil,
			},
		},
		{
			name:    "from db",
			builder: (&Selector[TestModel]{}).From("`test_db`.`test_model`"),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_db`.`test_model`;",
				Args: nil,
			},
		},
		{
			name:    "empty where",
			builder: (&Selector[TestModel]{}).Where(),
			wantQuery: &Query{
				SQL: "SELECT * FROM `TestModel`;",
				//Args: []any{123},
			},
		},
		{
			name:    "where",
			builder: (&Selector[TestModel]{}).Where(C("Age").Eq(123)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel` WHERE (`Age` = ?);",
				Args: []any{123},
			},
		},
		{
			name:    "not",
			builder: (&Selector[TestModel]{}).Where(Not(C("Age").Eq(123))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel` WHERE (NOT (`Age` = ?));",
				Args: []any{123},
			},
		},
		{
			name:    "not",
			builder: (&Selector[TestModel]{}).Where(Not(C("Age").Eq(123)).And(C("Id").Eq(321))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel` WHERE ((NOT (`Age` = ?))AND (`Id` = ?));",
				Args: []any{123, 321},
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
