package orm

import (
	"context"
	"errors"
	"strings"
)

// Selectable 标记接口，代表查找的是查找的列或聚合部分，SELECT XXX
type Selectable interface {
	selectable()
}

type Selector[T any] struct {
	table string
	//model *model.Model
	where []Predicate
	//sb    *strings.Builder
	//args  []any
	builder

	//cols []string
	cols []Selectable

	//db *DB
	//r *registry

	sess Session
}

//func (db *DB) NewSelector[T any]() *Selector[T] {
//	return &Selector[T]{
//		sb: &strings.Builder{},
//		db: db,
//	}
//}

func NewSelector[T any](sess Session) *Selector[T] {
	core := sess.getCore()
	return &Selector[T]{
		builder: builder{sb: &strings.Builder{}, core: core, quoter: core.dialect.quoter()},
		// dialect: sess.dialect, quoter: db.dialect.quoter(),

		//sb: &strings.Builder{},
		//db: db,
		sess: sess,
	}
}

//func (s *Selector[T]) Demo[S any]() (*Query, error) {
//
//}

func (s *Selector[T]) Build() (*Query, error) {
	//var sb strings.Builder
	//sb := s.sb
	//if s.sb == nil {
	//	s.sb = &strings.Builder{}
	//}
	var err error
	//r := &registry{}
	s.model, err = s.r.Registry(new(T))
	if err != nil {
		return nil, err
	}
	s.sb.WriteString("SELECT")

	if err = s.buildColumns(); err != nil {
		return nil, err
	}

	s.sb.WriteString("FROM ")
	// 反射拿到表名
	if s.table == "" {
		//var t T
		//Typ := reflect.TypeOf(t)
		//s.table = Typ.Name()
		s.sb.WriteByte('`')
		//s.sb.WriteString(Typ.Name())
		s.sb.WriteString(s.model.TableName)
		s.sb.WriteByte('`')
	} else {
		//sb.WriteByte('`')
		s.sb.WriteString(s.table)
		//sb.WriteByte('`')
	}

	// WHERE
	if len(s.where) > 0 {
		s.sb.WriteString(" WHERE ")
		//p := s.where[0]
		//for i := 1; i < len(s.where); i++ {
		//	p = p.And(s.where[i])
		//}
		//
		//if err := s.buildExpression(p); err != nil {
		//	return nil, err
		//}
		if err = s.buildPredicates(s.where); err != nil {
			return nil, err
		}
	}

	s.sb.WriteString(";")
	return &Query{
		SQL:  s.sb.String(),
		Args: s.args,
	}, nil
}

//func (s *Selector[T]) buildExpression(expr Expression) error {
//	switch exp := expr.(type) {
//	case nil:
//		return nil
//	case Predicate:
//		// 在这里构建p
//		// p.left构建好
//		// p.Op
//		// p.right
//		s.sb.WriteByte('(')
//		if err := s.buildExpression(exp.left); err != nil {
//			return err
//		}
//		//s.sb.WriteByte(' ')
//		if exp.op.String() != "" {
//			s.sb.WriteString(exp.op.String())
//			s.sb.WriteByte(' ')
//		}
//
//		if err := s.buildExpression(exp.right); err != nil {
//			return err
//		}
//		s.sb.WriteByte(')')
//	case Column:
//		//s.sb.WriteByte('`')
//		//fd, exist := s.model.FieldMap[exp.name]
//		//if !exist {
//		//	return errs.NewErrUnknownField(exp.name)
//		//}
//		//s.sb.WriteString(fd.ColName)
//		//s.sb.WriteByte('`')
//		//s.sb.WriteByte(' ')
//		exp.alias = ""
//		return s.buildColumn(exp)
//	case Op:
//	case value:
//		s.addArg(exp.val)
//		//s.args = append(s.args, raw.(value).val)
//		s.sb.WriteByte('?')
//	case RawExpr:
//		s.sb.WriteByte('(')
//		s.sb.WriteString(exp.raw)
//		s.addArg(exp.args...)
//		s.sb.WriteByte(')')
//	default:
//		return errs.NewErrUnsupportedExpression(expr)
//	}
//	return nil
//}

func (s *Selector[T]) buildColumns() (err error) {
	if len(s.cols) > 0 {
		//s.sb.WriteByte(' ')
		//for i, col := range s.cols {
		//	s.sb.WriteByte('`')
		//	s.sb.WriteString(col)
		//	s.sb.WriteByte('`')
		//	if i < len(s.cols)-1 {
		//		s.sb.WriteByte(',')
		//	}
		//	s.sb.WriteByte(' ')
		//}
		for i, col := range s.cols {
			s.sb.WriteByte(' ')
			switch exp := col.(type) {
			case Column:
				err = s.buildColumn(exp)
				if err != nil {
					return err
				}
			case Aggregate:
				s.sb.WriteString(exp.fn)
				s.sb.WriteByte('(')
				err = s.buildColumn(
					Column{
						name: exp.arg,
					})
				if err != nil {
					return err
				}
				s.sb.WriteByte(')')
				if exp.alias != "" {
					s.sb.WriteByte(' ')
					s.sb.WriteString("AS")
					s.sb.WriteByte(' ')
					s.sb.WriteByte('`')
					s.sb.WriteString(exp.alias)
					s.sb.WriteByte('`')
				}
			case RawExpr:
				s.sb.WriteString(exp.raw)
				s.addArg(exp.args...)
			default:
				return errors.New("")
			}

			if i < len(s.cols)-1 {
				s.sb.WriteByte(',')
			}
			//s.sb.WriteByte(' ')
		}
		s.sb.WriteByte(' ')
	} else {
		s.sb.WriteString(" * ")
	}
	return nil
}

//func (s *Selector[T]) buildColumn(column Column) error {
//	s.sb.WriteByte('`')
//	fd, exist := s.model.FieldMap[column.name]
//	if !exist {
//		return errs.NewErrUnknownField(column.name)
//	}
//	s.sb.WriteString(fd.ColName)
//	s.sb.WriteByte('`')
//
//	if column.alias != "" {
//		s.sb.WriteByte(' ')
//		s.sb.WriteString("AS")
//		s.sb.WriteByte(' ')
//		s.sb.WriteByte('`')
//		s.sb.WriteString(column.alias)
//		s.sb.WriteByte('`')
//	}
//
//	//s.sb.WriteByte(' ')
//	return nil
//}
//
//func (s *Selector[T]) addArg(vals ...any) *Selector[T] {
//	if len(vals) == 0 {
//		return nil
//	}
//	if s.args == nil {
//		s.args = make([]any, 0, 4) // 给定预估容量，避免频繁扩容
//	}
//	s.args = append(s.args, vals...)
//	return s
//}

// select最简实现
//func (s *Selector[T]) Select(cols ...string) *Selector[T] {
//	s.cols = cols
//	return s
//}
//
//func (s *Selector[T]) Select(col string) *Selector[T] {
//	s.col = col
//	return s
//}

func (s *Selector[T]) Select(cols ...Selectable) *Selector[T] {
	s.cols = cols
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

// 基于reflect
//func (s Selector[T]) Get(ctx context.Context) (*T, error) {
//	q, err := s.Build()
//	if err != nil {
//		return nil, err
//	}
//
//	//var db *sql.DB
//	db := s.db.db
//	// 发起查询，处理结果集
//	row, err := db.QueryContext(ctx, q.SQL, q.Args...)
//	if err != nil {
//		return nil, err
//	}
//
//	if !row.Next() {
//		// 里面是否返回error，返回error和sql包一致吗？和GetMulti保持一致
//		return nil, ErrNoRows
//	}
//
//	//s.model.FieldMap
//
//	// 问题： 类型、顺序要匹配
//
//	// select出来哪些列
//	cs, err := row.Columns()
//	if err != nil {
//		return nil, err
//	}
//
//	tp := new(T)
//	vals := make([]any, 0, len(cs))
//	valElems := make([]reflect.Value, 0, len(cs))
//	for _, c := range cs {
//		fd, ok := s.model.ColumnMap[c]
//		if !ok {
//			return nil, errs.NewErrUnknownColumn(c)
//		}
//		val := reflect.New(fd.Typ)
//		vals = append(vals, val.Interface())
//		valElems = append(valElems, val.Elem())
//
//		//for _, fd := range s.model.FieldMap {
//		//	if fd.ColName == c {
//		//		// 反射创建新的实例
//		//		val := reflect.New(fd.Typ)
//		//		vals = append(vals, val.Interface())
//		//	}
//		//}
//	}
//	err = row.Scan(vals...)
//	if err != nil {
//		return nil, err
//	}
//
//	//tpValue := reflect.ValueOf(tp)
//	//for i, c := range cs {
//	//	fd, ok := s.model.ColumnMap[c]
//	//	if !ok {
//	//		return nil, errs.NewErrUnknownColumn(c)
//	//	}
//	//	tpValue.Elem().FieldByName(fd.GoName).Set(valElems[i])
//	//	//tpValue.Elem().FieldByName(fd.GoName).Set(reflect.ValueOf(vals[i]).Elem())
//	//	//for _, fd := range s.model.FieldMap {
//	//	//	if fd.ColName == c {
//	//	//		tpValue.Elem().FieldByName(fd.GoName).Set(reflect.ValueOf(vals[i]).Elem())
//	//	//	}
//	//	//}
//	//}
//
//	tpValueElem := reflect.ValueOf(tp).Elem()
//	for i, c := range cs {
//		fd, ok := s.model.ColumnMap[c]
//		if !ok {
//			return nil, errs.NewErrUnknownColumn(c)
//		}
//		tpValueElem.FieldByName(fd.GoName).Set(valElems[i])
//		//tpValue.Elem().FieldByName(fd.GoName).Set(reflect.ValueOf(vals[i]).Elem())
//		//for _, fd := range s.model.FieldMap {
//		//	if fd.ColName == c {
//		//		tpValue.Elem().FieldByName(fd.GoName).Set(reflect.ValueOf(vals[i]).Elem())
//		//	}
//		//}
//	}
//
//	return tp, nil
//
//}

// 基于unsafe
func (s Selector[T]) Get(ctx context.Context) (*T, error) {
	q, err := s.Build()
	if err != nil {
		return nil, err
	}

	//var db *sql.DB
	sess := s.sess
	// 发起查询，处理结果集
	row, err := sess.queryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return nil, err
	}

	if !row.Next() {
		// 里面是否返回error，返回error和sql包一致吗？和GetMulti保持一致
		return nil, ErrNoRows
	}

	tp := new(T)
	//var creator valuer.Creator
	err = s.creator(s.model, tp).SetColumns(row)
	return tp, err

	////s.model.FieldMap
	//
	//// 问题： 类型、顺序要匹配
	//
	//// select出来哪些列
	//cs, err := row.Columns()
	//if err != nil {
	//	return nil, err
	//}
	//
	//tp := new(T)
	//vals := make([]any, 0, len(cs))
	//address := reflect.ValueOf(tp).UnsafePointer()
	//for _, c := range cs {
	//	fd, ok := s.model.ColumnMap[c]
	//	if !ok {
	//		return nil, errs.NewErrUnknownColumn(c)
	//	}
	//	// 起始地址+偏移量
	//	fdAddress := unsafe.Pointer(uintptr(address) + fd.Offset)
	//
	//	// 反射在特定地址上，创建特定类型实例，原本类型的指针类型；case：fd.Typ=int, val是*int
	//	val := reflect.NewAt(fd.Typ, fdAddress)
	//	vals = append(vals, val.Interface())
	//}
	//err = row.Scan(vals...)
	//if err != nil {
	//	return nil, err
	//}

	return tp, nil

}

func (s Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	q, err := s.Build()
	if err != nil {
		return nil, err
	}

	//var db *sql.DB
	sess := s.sess
	// 发起查询，处理结果集
	rows, err := sess.queryContext(ctx, q.SQL, q.Args)

	for rows.Next() {

	}
	return nil, nil
}
