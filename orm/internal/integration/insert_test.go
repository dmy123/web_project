//go:build e2e

package integration

import (
	"awesomeProject1/orm"
	"awesomeProject1/orm/internal/test"
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type InsertTestSuite struct {
	Suite
}

func TestMySQLInsert(t *testing.T) {
	suite.Run(t, &InsertTestSuite{Suite: Suite{
		driver: "mysql",
		dsn:    "root:root@tcp(localhost:13306)/integration_test",
	},
	})
}

func (i *InsertTestSuite) TearDownTest() {
	res := orm.RawQuery[any](i.db, "TRUNCATE TABLE `simple_struct`").
		Exec(context.Background())
	require.NoError(i.T(), res.Err())
}

func (i *InsertTestSuite) TestInsert() {
	db := i.db
	t := i.T()
	testCase := []struct {
		name     string
		i        *orm.Inserter[test.SimpleStruct]
		affected int64
	}{
		{
			name:     "insert one",
			i:        orm.NewInserter[test.SimpleStruct](db).Values(test.NewSimpleStruct(37)),
			affected: 1,
		},
		{
			name:     "insert multiple",
			i:        orm.NewInserter[test.SimpleStruct](db).Values(test.NewSimpleStruct(38), test.NewSimpleStruct(39)),
			affected: 2,
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			_, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			res := tc.i.Exec(context.Background())
			affected, err := res.RowsAffected()
			assert.NoError(t, res.Err())
			assert.NoError(t, err)
			assert.Equal(t, affected, tc.affected)
		})
	}
}

//func TestInsert(t *testing.T) {
//	db, err := orm.Open("mysql", "root:root@tcp(localhost:13306)/integration_test")
//	assert.NoError(t, err)
//	testCase := []struct {
//		name     string
//		i        *orm.Inserter[test.SimpleStruct]
//		affected int64
//	}{
//		{
//			name:     "insert one",
//			i:        orm.NewInserter[test.SimpleStruct](db).Values(test.NewSimpleStruct(17)),
//			affected: 1,
//		},
//		{
//			name:     "insert multiple",
//			i:        orm.NewInserter[test.SimpleStruct](db).Values(test.NewSimpleStruct(18), test.NewSimpleStruct(19)),
//			affected: 2,
//		},
//	}
//	for _, tc := range testCase {
//		t.Run(tc.name, func(t *testing.T) {
//			_, cancel := context.WithTimeout(context.Background(), time.Second*10)
//			defer cancel()
//			res := tc.i.Exec(context.Background())
//			affected, err := res.RowsAffected()
//			assert.NoError(t, res.Err())
//			assert.NoError(t, err)
//			assert.Equal(t, affected, tc.affected)
//		})
//	}
//}

//type SQLite3InsertTest struct {
//	InsertTestSuite
//}
//
//func (i *SQLite3InsertTest) SetupSuite() {
//	db, err := sql.Open(i.driver, i.dsn)
//	// 建表
//	db.ExecContext(context.Background(), "")
//	require.NoError(i.T(), err)
//	i.db, err = orm.OpenDB(db)
//	require.NoError(i.T(), err)
//}
//
//func TestSQLite3(t *testing.T) {
//	suite.Run(t, &SQLite3InsertTest{InsertTestSuite{Suite: Suite{
//		driver: "sqlite3",
//		dsn:    "file:test.db?cache=shared&mode=memory",
//	},
//	},
//	})
//}
