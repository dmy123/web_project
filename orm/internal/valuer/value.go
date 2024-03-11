package valuer

import (
	"awesomeProject1/orm/model"
	"database/sql"
)

type Value interface {
	Field(name string) (any, error)
	SetColumns(rows *sql.Rows) error
}

// 简单的函数式工厂接口
type Creator func(model *model.Model, entity any) Value
