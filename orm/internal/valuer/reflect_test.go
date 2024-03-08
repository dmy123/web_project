package valuer

import (
	"awesomeProject1/orm/model"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_reflectValue_SetColumns(t *testing.T) {
	testSetColumns(t, NewReflectValue)
}

func testSetColumns(t *testing.T, creator Creator) {
	type fields struct {
		// 一定是指针
		val any
	}
	type args struct {
		rows *sqlmock.Rows
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    error
		wantEntity any
	}{
		// TODO: Add test cases.
		{
			name: "set column",
			fields: fields{
				val: &TestModel{},
			},
			args: args{
				rows: func() *sqlmock.Rows {
					rows := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
					rows.AddRow("1", "tom", 18, "jerry")
					return rows
				}(),
			},
			wantEntity: &TestModel{
				Id:        1,
				FirstName: "tom",
				Age:       18,
				LastName: &sql.NullString{
					String: "jerry",
					Valid:  true,
				},
			},
		},
		{
			name: "order",
			fields: fields{
				val: &TestModel{},
			},
			args: args{
				rows: func() *sqlmock.Rows {
					rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age"})
					rows.AddRow("1", "tom", "jerry", 18)
					return rows
				}(),
			},
			wantEntity: &TestModel{
				Id:        1,
				FirstName: "tom",
				Age:       18,
				LastName: &sql.NullString{
					String: "jerry",
					Valid:  true,
				},
			},
		},
		{
			name: "partial columns",
			fields: fields{
				val: &TestModel{},
			},
			args: args{
				rows: func() *sqlmock.Rows {
					rows := sqlmock.NewRows([]string{"id", "last_name"})
					rows.AddRow("1", "jerry")
					return rows
				}(),
			},
			wantEntity: &TestModel{
				Id: 1,
				//FirstName: "tom",
				//Age:       18,
				LastName: &sql.NullString{
					String: "jerry",
					Valid:  true,
				},
			},
		},
	}
	r := model.MustNewRegistry()

	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 构造rows
			mockRows := tt.args.rows
			mock.ExpectQuery("SELECT XX").WillReturnRows(mockRows)
			rows, err := mockDB.Query("SELECT XX")
			assert.NoError(t, err)

			rows.Next()

			model, err := r.Get(tt.fields.val)
			assert.NoError(t, err)
			//ref := NewReflectValue(model, tt.fields.val)

			//err = ref.SetColumns(rows)
			err = creator(model, tt.fields.val).SetColumns(rows)
			assert.Equal(t, err, tt.wantErr)
			if err != nil {
				return
			}
			assert.Equal(t, tt.wantEntity, tt.fields.val)
		})
	}
}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}
