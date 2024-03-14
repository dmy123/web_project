package orm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRawQuerier_Get(t *testing.T) {
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
		r    *RawQuerier[T]
		//args    args
		want    *T
		wantErr error
	}
	tests := []testCase[TestModel]{
		{
			name:    "query error",
			r:       RawQuery[TestModel](db, "SELECT * FROM `test_model`"),
			wantErr: errors.New("query error"),
		},
		{
			name:    "no rows",
			r:       RawQuery[TestModel](db, "SELECT * FROM `test_model` WHERE `id` = -1", -1),
			wantErr: ErrNoRows,
		},
		{
			name: "data",
			r:    RawQuery[TestModel](db, "SELECT * FROM `test_model` WHERE `id` = ?", 1),
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
		//	r:       NewSelector[TestModel](db).Where(C("Id").Lt(1)),
		//	wantErr: errs.ErrNoRows,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.Get(context.Background())
			if err != nil || tt.wantErr != nil {
				assert.Equal(t, err, tt.wantErr)
				return
			}
			assert.Equalf(t, tt.want, got, "Get(%v)", context.Background())
		})
	}
}
