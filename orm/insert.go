package orm

import (
	"awesomeProject1/orm/internal/errs"
	"awesomeProject1/orm/model"
	"context"
	"strings"
)

type UpsertBuilder[T any] struct {
	i               *Inserter[T]
	conflictColumns []string
}

type Upsert struct {
	assigns         []Assignable
	conflictColumns []string
}

func (o *UpsertBuilder[T]) ConflictColumns(cols ...string) *UpsertBuilder[T] {
	o.conflictColumns = cols
	return o
}

func (o *UpsertBuilder[T]) Update(assigns ...Assignable) *Inserter[T] {
	o.i.onDuplicateKey = &Upsert{
		assigns:         assigns,
		conflictColumns: o.conflictColumns,
	}
	return o.i
}

func (i *Inserter[T]) OnDuplicateKey() *UpsertBuilder[T] {
	return &UpsertBuilder[T]{
		i: i,
	}
}

type Assignable interface {
	assign()
}

type Inserter[T any] struct {
	//db   *DB
	sess Session
	vals []*T
	builder
	//sb      *strings.Builder
	//model   *model.Model
	//args    []any
	columns []string
	//onDuplicateKey []Assignable
	onDuplicateKey *Upsert
}

func NewInserter[T any](sess Session) *Inserter[T] {
	core := sess.getCore()
	return &Inserter[T]{
		builder: builder{sb: &strings.Builder{}, core: core, quoter: core.dialect.quoter()},
		sess:    sess,
	}
}

//// 对非mysql的库不友好
//func (i *Inserter[T]) Upsert(assigns ...Assignable) *Inserter[T] {
//	i.onDuplicateKey = assigns
//	return i
//}

// Values 指定插入哪些数据
func (i *Inserter[T]) Values(vals ...*T) *Inserter[T] {
	i.vals = vals
	return i
}

func (i *Inserter[T]) Columns(cols ...string) *Inserter[T] {
	i.columns = cols
	return i
}

func (i *Inserter[T]) Build() (res *Query, err error) {
	if len(i.vals) == 0 {
		return nil, errs.ErrInsertZeroRows
	}
	i.sb.WriteString("INSERT INTO ")
	i.model, err = i.r.Get(new(T))
	if err != nil {
		return nil, err
	}
	i.sb.WriteByte('`')
	i.sb.WriteString(i.model.TableName)
	i.sb.WriteByte('`')
	i.sb.WriteByte(' ')
	i.sb.WriteByte('(')

	fields := i.model.Fields
	if len(i.columns) > 0 {
		fields = make([]*model.Field, 0, len(i.columns))
		for _, fd := range i.columns {
			f, ok := i.model.FieldMap[fd]
			if !ok {
				return nil, errs.NewErrUnknownField(fd)
			}
			fields = append(fields, f)
		}
	}

	// 显示指定列的顺序，不然不知道数据库默认顺序
	cnt := 0
	//for _, fd := range i.model.FieldMap {
	//	if cnt > 0 {
	//		i.sb.WriteString(", ")
	//	}
	//	i.sb.WriteByte('`')
	//	i.sb.WriteString(fd.ColName)
	//	i.sb.WriteByte('`')
	//	cnt++
	//}
	for _, fd := range fields {
		if cnt > 0 {
			i.sb.WriteString(", ")
		}
		i.Quoter(fd.ColName)
		//i.sb.WriteByte('`')
		//i.sb.WriteString(fd.ColName)
		//i.sb.WriteByte('`')
		cnt++
	}

	i.sb.WriteByte(')')
	i.sb.WriteString(" VALUES ")
	i.args = make([]any, 0, len(i.vals)*len(fields))
	for k, val := range i.vals {
		if k > 0 {
			i.sb.WriteByte(',')
		}
		i.sb.WriteByte('(')
		v := i.creator(i.model, val)
		for j := 0; j < cnt; j++ {
			if j > 0 {
				i.sb.WriteByte(',')
			}
			i.sb.WriteByte('?')
			arg, err := v.Field(fields[j].GoName)
			if err != nil {
				return nil, err
			}
			//v := reflect.ValueOf(val).Elem().FieldByName(fields[j].GoName).Interface()
			//i.args = append(i.args, v)
			//i.addArg(v)
			i.addArg(arg)
		}
		i.sb.WriteByte(')')
	}

	if i.onDuplicateKey != nil {
		err = i.dialect.buildUpsert(&i.builder, i.onDuplicateKey)
		if err != nil {
			return nil, err
		}
		//i.sb.WriteString(" ON DUPLICATE KEY UPDATE ")
		//for j, assign := range i.onDuplicateKey.assigns {
		//	switch exp := assign.(type) {
		//	case Assignment:
		//		fd, ok := i.model.FieldMap[exp.col]
		//		if !ok {
		//			return nil, errs.NewErrUnknownField(exp.col)
		//		}
		//		if j > 0 {
		//			i.sb.WriteString(", ")
		//		}
		//		i.sb.WriteByte('`')
		//		i.sb.WriteString(fd.ColName)
		//		i.sb.WriteByte('`')
		//		i.sb.WriteString("=?")
		//		i.args = append(i.args, exp.val...)
		//	case Column:
		//		fd, ok := i.model.FieldMap[exp.name]
		//		if !ok {
		//			return nil, errs.NewErrUnknownField(exp.name)
		//		}
		//		if j > 0 {
		//			i.sb.WriteString(", ")
		//		}
		//		i.sb.WriteByte('`')
		//		i.sb.WriteString(fd.ColName)
		//		i.sb.WriteByte('`')
		//		i.sb.WriteString("=")
		//		i.sb.WriteString("VALUES")
		//		i.sb.WriteByte('(')
		//		i.sb.WriteByte('`')
		//		i.sb.WriteString(fd.ColName)
		//		i.sb.WriteByte('`')
		//		i.sb.WriteByte(')')
		//	default:
		//		return nil, errs.NewErrUnsupportedAssignable(exp)
		//	}
		//}
	}

	i.sb.WriteByte(';')

	return &Query{
		SQL:  i.sb.String(),
		Args: i.args,
	}, nil
}

func (i *Inserter[T]) Exec(ctx context.Context) Result {
	q, err := i.Build()
	if err != nil {
		return Result{
			err: err,
		}
	}
	r, err := i.sess.execContext(ctx, q.SQL, q.Args...)
	return Result{
		res: r,
		err: err,
	}
}
