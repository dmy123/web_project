package sql

import (
	"context"
	"errors"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestSQLMock(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer db.Close() // 注意close
	require.NoError(t, err)

	mockRows := sqlmock.NewRows([]string{"id", "first_name"})
	mockRows.AddRow(1, "Tom")
	// 正则表达式
	// mock顺序需要与查询顺序一致，否则失败。
	mock.ExpectQuery("SELECT id,first_name FROM `user`.*").WillReturnRows(mockRows)
	mock.ExpectQuery("SELECT id FROM `user`.*").WillReturnError(errors.New("mock error"))
	mock.ExpectQuery("SELECT id,first_name FROM `user`.*").WillReturnRows(mockRows)

	//
	//result := sqlmock.NewResult()

	rows, err := db.QueryContext(context.Background(), "SELECT id,first_name FROM `user`;")
	require.NoError(t, err)
	for rows.Next() {
		tm := TestModel{}
		err = rows.Scan(&tm.Id, &tm.FirstName)
		require.NoError(t, err)
		log.Println("数据", tm)
	}

	_, err = db.QueryContext(context.Background(), "SELECT id FROM `user`;")
	require.Error(t, err)
}
