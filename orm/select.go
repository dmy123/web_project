package orm

import (
	"context"
	"reflect"
	"strings"
)

type Selector[T any] struct {
	table string
}

func (s Selector[T]) Build() (*Query, error) {
	var sb strings.Builder
	sb.WriteString("SELECT * FROM ")
	// 反射拿到表名
	if s.table == "" {
		var t T
		typ := reflect.TypeOf(t)
		//s.table = typ.Name()
		sb.WriteByte('`')
		sb.WriteString(typ.Name())
		sb.WriteByte('`')
	} else {
		//sb.WriteByte('`')
		sb.WriteString(s.table)
		//sb.WriteByte('`')
	}
	sb.WriteString(";")
	return &Query{
		SQL: sb.String(),
	}, nil
}

func (s *Selector[T]) From(table string) *Selector[T] {
	s.table = table
	return s
}

func (s Selector[T]) Get(ctx context.Context) (*T, error) {
	//TODO implement me
	panic("implement me")
}

func (s Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	//TODO implement me
	panic("implement me")
}
