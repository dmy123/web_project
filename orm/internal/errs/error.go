package errs

import (
	"errors"
	"fmt"
)

var (
	ErrPointerOnly = errors.New("orm：只支持指向结构体一级指针")
)

func NewErrUnsupportedExpression(expr any) error {
	return fmt.Errorf("orm：不支持的表达式类型 %v", expr)
}

func NewErrUnknownField(name string) error {
	return fmt.Errorf("orm：未知字段 %s", name)
}

func NewErrInvalidTagContent(pair string) error {
	return fmt.Errorf("orm：非法标签值 %s", pair)
}
