package valuer

import "database/sql"

type Valuer interface {
	SetColumns(rows *sql.Rows) error
}

// 简单的函数式工厂接口
type Creator func(entity any) Valuer
