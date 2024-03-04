package orm

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

type Selector[T any] struct {
	table string
	where []Predicate
	sb    *strings.Builder
	args  []any
}

func (s Selector[T]) Build() (*Query, error) {
	//var sb strings.Builder
	//sb := s.sb
	if s.sb == nil {
		s.sb = &strings.Builder{}
	}
	s.sb.WriteString("SELECT * FROM ")
	// 反射拿到表名
	if s.table == "" {
		var t T
		typ := reflect.TypeOf(t)
		//s.table = typ.Name()
		s.sb.WriteByte('`')
		s.sb.WriteString(typ.Name())
		s.sb.WriteByte('`')
	} else {
		//sb.WriteByte('`')
		s.sb.WriteString(s.table)
		//sb.WriteByte('`')
	}

	// WHERE
	if len(s.where) > 0 {
		s.sb.WriteString(" WHERE ")
		p := s.where[0]
		for i := 1; i < len(s.where); i++ {
			p = p.And(s.where[i])
		}

		if err := s.buildExpression(p); err != nil {
			return nil, err
		}
	}

	s.sb.WriteString(";")
	return &Query{
		SQL:  s.sb.String(),
		Args: s.args,
	}, nil
}

func (s *Selector[T]) buildExpression(expr Expression) error {
	switch expr.(type) {
	case nil:
		return nil
	case Predicate:
		// 在这里构建p
		// p.left构建好
		// p.Op
		// p.right
		s.sb.WriteByte('(')
		if err := s.buildExpression(expr.(Predicate).left); err != nil {
			return err
		}
		//s.sb.WriteByte(' ')
		s.sb.WriteString(expr.(Predicate).op.String())
		s.sb.WriteByte(' ')

		if err := s.buildExpression(expr.(Predicate).right); err != nil {
			return err
		}
		s.sb.WriteByte(')')
	case Column:
		s.sb.WriteByte('`')
		s.sb.WriteString(expr.(Column).name)
		s.sb.WriteByte('`')
		s.sb.WriteByte(' ')
	case Op:
	case value:
		s.addArg(expr.(value).val)
		//s.args = append(s.args, expr.(value).val)
		s.sb.WriteByte('?')
	default:
		return fmt.Errorf("orm: 不支持的表达式类型 %v", expr)
	}
	return nil
}

func (s *Selector[T]) addArg(val any) *Selector[T] {
	if s.args == nil {
		s.args = make([]any, 0, 4) // 给定预估容量，避免频繁扩容
	}
	s.args = append(s.args, val)
	return s
}

func (s *Selector[T]) From(table string) *Selector[T] {
	s.table = table
	return s
}

func (s *Selector[T]) Where(ps ...Predicate) *Selector[T] {
	s.where = ps
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
