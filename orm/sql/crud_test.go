package sql

import (
	"context"
	"database/sql"
	//_ "github.com/mattn/go-sqlite3"
	_ "github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func TestDB(t *testing.T) {
	db, err := sql.Open("sqlite", "file:test.db?cache=shared&mode=memory")
	require.NoError(t, err)
	defer db.Close()
	db.Ping()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// select语句用query，其他用ExecContext
	_, err = db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS test_model(
    id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    age INTEGER,
    last_name TEXT NOT NULL
)
`)

	require.NoError(t, err)

	// 用？当占位符，防止sql注入
	res, err := db.ExecContext(ctx, "INSERT INTO `test_model`(`id`, `first_name`, `age`, `last_name`) VALUES(?,?,?,?)", 1, "tom", 18, "jerry")

	require.NoError(t, err)
	affected, err := res.RowsAffected()
	require.NoError(t, err)
	log.Println("受影响行 ", affected)
	lastId, err := res.LastInsertId()
	require.NoError(t, err)
	log.Println("最后输入 ", lastId)

	row := db.QueryRowContext(ctx, "SELECT * FROM `test_model` WHERE `id` = ?", 1)
	require.NoError(t, row.Err())
	tm := TestModel{}
	err = row.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
	require.NoError(t, err)
	log.Println("1、数据", tm)

	row = db.QueryRowContext(ctx, "SELECT * FROM `test_model` WHERE `id` = ?", 2)
	require.NoError(t, row.Err())
	log.Println("2、err", row.Err())

	row = db.QueryRowContext(ctx, "SELECT * FROM `test_model` WHERE `id` = ?", 2)
	tm = TestModel{}
	err = row.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
	require.Error(t, sql.ErrNoRows, err)
	log.Println("3、数据", tm)
	log.Println("3、err", err.Error())

	rows, err := db.QueryContext(ctx, "SELECT * FROM `test_model` WHERE `id` = ?", 1)
	require.NoError(t, err)
	require.NoError(t, rows.Err())
	for rows.Next() {
		tm = TestModel{}
		err = rows.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
		require.NoError(t, err)
		log.Println("4、数据", tm)
	}

	rows, err = db.QueryContext(ctx, "SELECT * FROM `test_model` WHERE `id` = ?", 2)
	require.NoError(t, err)
	require.NoError(t, rows.Err())
	for rows.Next() {
		tm = TestModel{}
		err = rows.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
		require.NoError(t, err)
		log.Println("5、数据", tm)
	}

	// Scan才会返回查询无数据的error

	cancel()
}

func TestTX(t *testing.T) {
	db, err := sql.Open("sqlite", "file:test.db?cache=shared&mode=memory")
	require.NoError(t, err)
	defer db.Close()
	db.Ping()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// select语句用query，其他用ExecContext
	_, err = db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS test_model(
    id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    age INTEGER,
    last_name TEXT NOT NULL
)
`)
	require.NoError(t, err)

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})

	// 用？当占位符，防止sql注入
	res, err := tx.ExecContext(ctx, "INSERT INTO `test_model`(`id`, `first_name`, `age`, `last_name`) VALUES(?,?,?,?)", 1, "tom", 18, "jerry")
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			log.Println(err)
		}
		return
	}

	require.NoError(t, err)
	affected, err := res.RowsAffected()
	require.NoError(t, err)
	log.Println("受影响行 ", affected)
	lastId, err := res.LastInsertId()
	require.NoError(t, err)
	log.Println("最后输入 ", lastId)

	require.NoError(t, err)
	err = tx.Commit()
	require.NoError(t, err)

	db.Close()

	cancel()
}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}
