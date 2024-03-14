package errs

import (
	"errors"
	"fmt"
)

var (
	ErrPointerOnly    = errors.New("orm：只支持指向结构体一级指针")
	ErrNoRows         = errors.New("orm: 没有数据")
	ErrInsertZeroRows = errors.New("orm: 插入0行")
)

func NewErrUnsupportedExpression(expr any) error {
	return fmt.Errorf("orm：不支持的表达式类型 %v", expr)
}

func NewErrFailedToRollbackTx(bizErr error, rbErr error, panicked bool) error {
	return fmt.Errorf("orm：事务闭包回滚失败，业务错误 %w，回滚错误 %s， 是否panic %t", bizErr, rbErr, panicked)
}

func NewErrUnknownField(name string) error {
	return fmt.Errorf("orm：未知字段 %s", name)
}

func NewErrUnknownColumn(name string) error {
	return fmt.Errorf("orm：未知列 %s", name)
}

func NewErrInvalidTagContent(pair string) error {
	return fmt.Errorf("orm：非法标签值 %s", pair)
}

func NewErrUnsupportedAssignable(expr any) error {
	return fmt.Errorf("orm：不支持的表达式类型 %v", expr)
}

func NewErrUnsupportedTable(table any) error {
	return fmt.Errorf("orm：不支持的TableReference类型 %v", table)
}
