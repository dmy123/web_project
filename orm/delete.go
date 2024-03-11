package orm

import (
	"strings"
)

type Deleter[T any] struct {
	//sb    *strings.Builder
	//model *model.Model
	db *DB

	wheres []Predicate
	//args   []any
	table string
	builder
}

func NewDeleter[T any](db *DB) *Deleter[T] {
	return &Deleter[T]{builder: builder{sb: &strings.Builder{}, dialect: db.dialect, quoter: db.dialect.quoter()},
		db: db}
}

func (d Deleter[T]) Build() (query *Query, err error) {
	d.model, err = d.db.r.Registry(new(T))
	// 构造语句
	d.sb.WriteString("DELETE FROM ")

	if len(d.table) == 0 {
		d.sb.WriteByte('`')
		d.sb.WriteString(d.model.TableName)
		d.sb.WriteByte('`')
	} else {
		d.sb.WriteString(d.table)
	}

	// 处理where
	if len(d.wheres) > 0 {
		d.sb.WriteString(" WHERE ")
		if err = d.buildPredicates(d.wheres); err != nil {
			return nil, err
		}
	}

	d.sb.WriteByte(';')
	return &Query{
		SQL:  d.sb.String(),
		Args: d.args,
	}, nil
}

func (d *Deleter[T]) From(table string) *Deleter[T] {
	d.table = table
	return d
}

func (d *Deleter[T]) Where(predicates ...Predicate) *Deleter[T] {
	d.wheres = predicates
	return d
}

//func (d *Deleter[T]) buildExpression(expr Expression) error {
//	switch exp := expr.(type) {
//	case Predicate:
//		d.buildExpression(exp.left)
//		if len(exp.op) > 0 {
//			d.sb.WriteByte(' ')
//			d.sb.WriteString(exp.op.String())
//			d.sb.WriteByte(' ')
//		}
//		d.buildExpression(exp.right)
//	case Column:
//		d.sb.WriteByte('`')
//		d.sb.WriteString(exp.name)
//		d.sb.WriteByte('`')
//	case value:
//		d.sb.WriteByte('?')
//		d.args = append(d.args, exp.val)
//	default:
//		return errors.New("")
//	}
//	return nil
//}
