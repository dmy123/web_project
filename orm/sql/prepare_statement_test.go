package sql

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func TestPrepareStatement(t *testing.T) {
	db, err := sql.Open("sqlite", "file:test.db?cache=shared&mode=memory")
	require.NoError(t, err)
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	stmt, err := db.PrepareContext(ctx, "SELECT * FROM `test_model` WHERE `id`=?;")
	require.NoError(t, err)
	rows, err := stmt.QueryContext(ctx, 1)
	//require.NoError(t, err)

	//require.NoError(t, rows.Err())
	if err == nil && rows.Err() == nil {
		for rows.Next() {
			tm := TestModel{}
			err = rows.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
			require.NoError(t, err)
			log.Println("数据", tm)
		}
	}

	cancel()
	stmt.Close() // 整个应用关闭时调用

	stmt, err = db.PrepareContext(ctx, "SELECT * FROM `test_model` WHERE `id` IN (?,?,?);")
	stmt, err = db.PrepareContext(ctx, "SELECT * FROM `test_model` WHERE `id` IN (?,?,?,?);")
}
