package orm

import (
	"awesomeProject1/orm/internal/errs"
	"strings"
)

type builder struct {
	sb   *strings.Builder
	args []any

	core
	//model *model.Model
	//
	//dialect Dialect
	quoter byte
}

func (b *builder) buildPredicates(ps []Predicate) error {
	p := ps[0]
	for i := 1; i < len(ps); i++ {
		p = p.And(ps[i])
	}

	if err := b.buildExpression(p); err != nil {
		return err
	}
	return nil
}

func (s *builder) buildExpression(expr Expression) error {
	switch exp := expr.(type) {
	case nil:
		return nil
	case Predicate:
		s.sb.WriteByte('(')
		if err := s.buildExpression(exp.left); err != nil {
			return err
		}
		//s.sb.WriteByte(' ')
		if exp.op.String() != "" {
			s.sb.WriteString(exp.op.String())
			s.sb.WriteByte(' ')
		}

		if err := s.buildExpression(exp.right); err != nil {
			return err
		}
		s.sb.WriteByte(')')
	case Column:
		exp.alias = ""
		return s.buildColumn(exp)
	case Op:
	case value:
		s.addArg(exp.val)
		s.sb.WriteByte('?')
	case RawExpr:
		s.sb.WriteByte('(')
		s.sb.WriteString(exp.raw)
		s.addArg(exp.args...)
		s.sb.WriteByte(')')
	default:
		return errs.NewErrUnsupportedExpression(expr)
	}
	return nil
}

func (s *builder) buildColumn(column Column) error {
	s.sb.WriteByte('`')
	fd, exist := s.model.FieldMap[column.name]
	if !exist {
		return errs.NewErrUnknownField(column.name)
	}
	s.sb.WriteString(fd.ColName)
	s.sb.WriteByte('`')

	if column.alias != "" {
		s.sb.WriteByte(' ')
		s.sb.WriteString("AS")
		s.sb.WriteByte(' ')
		s.sb.WriteByte('`')
		s.sb.WriteString(column.alias)
		s.sb.WriteByte('`')
	}

	return nil
}

func (s *builder) addArg(vals ...any) *builder {
	if len(vals) == 0 {
		return nil
	}
	if s.args == nil {
		s.args = make([]any, 0, 4) // 给定预估容量，避免频繁扩容
	}
	s.args = append(s.args, vals...)
	return s
}

func (b *builder) Quoter(name string) {
	b.sb.WriteByte(b.quoter)
	b.sb.WriteString(name)
	b.sb.WriteByte(b.quoter)
}
