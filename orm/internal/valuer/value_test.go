package valuer

import (
	model2 "awesomeProject1/orm/model"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func BenchmarkSetColumns(b *testing.B) {

	fn := func(b *testing.B, creator Creator) {
		mockDB, mock, err := sqlmock.New()
		require.NoError(b, err)
		defer mockDB.Close()

		mockRows := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
		row := []driver.Value{"1", "tom", 18, "jerry"}
		for i := 0; i < b.N; i++ {
			mockRows.AddRow(row...)
		}
		mock.ExpectQuery("SELECT XX").WillReturnRows(mockRows)
		rows, err := mockDB.Query("SELECT XX")

		model, err := model2.MustNewRegistry().Get(&TestModel{})
		assert.NoError(b, err)
		// 重置计时器
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			rows.Next()

			//value := NewReflectValue(model, &TestModel{})
			value := creator(model, &TestModel{})
			value.SetColumns(rows)
		}
	}

	b.Run("reflect", func(b *testing.B) {
		fn(b, NewReflectValue)
	})

	b.Run("unsafe", func(b *testing.B) {
		fn(b, NewUnsafeValue)
	})
}
