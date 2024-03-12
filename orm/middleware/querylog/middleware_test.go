package querylog

import (
	"awesomeProject1/orm"
	"context"
	"database/sql"
	_ "github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMiddlewareBuilder_Builder(t *testing.T) {
	var q string
	var args []any
	m := (&MiddlewareBuilder{}).LogFunc(func(query string, as []any) {
		q = query
		args = as
	})

	db, err := orm.Open("sqlite", "file:test.db?cache=shared&mode=memory", orm.DBWithDialect(orm.DialectMySQL),
		orm.DBWithMiddleware(m.Builder()))
	require.NoError(t, err)

	_, _ = orm.NewSelector[TestModel](db).Where(orm.C("Id").Eq(10)).Get(context.Background())
	assert.Equal(t, "SELECT * FROM `test_model` WHERE (`id`= ?);", q)
	assert.Equal(t, []any{10}, args)

	orm.NewInserter[TestModel](db).Values(&TestModel{Id: 18}).Exec(context.Background())
	assert.Equal(t, "INSERT INTO `test_model` (`id`, `first_name`, `age`, `last_name`) VALUES (?,?,?,?);", q)
	assert.Equal(t, []any{int64(18), "", int8(0), (*sql.NullString)(nil)}, args)

	orm.NewDeleter[TestModel](db).Exec(context.Background())
	assert.Equal(t, "DELETE FROM `test_model`;", q)
	//var temp any
	//assert.Equal(t, []any{temp}, args)
}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}
