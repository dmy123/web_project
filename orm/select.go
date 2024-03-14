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

//func (r *Selector[T]) Demo[S any]() (*Query, error) {
//
//}

func (s *Selector[T]) Build() (*Query, error) {
	//var sb strings.Builder
	//sb := r.sb
	//if r.sb == nil {
	//	r.sb = &strings.Builder{}
	//}
	var err error
	//r := &registry{}
	if s.model == nil {
		s.model, err = s.r.Get(new(T))
		if err != nil {
			return nil, err
		}
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
		//r.table = Typ.Name()
		s.sb.WriteByte('`')
		//r.sb.WriteString(Typ.Name())
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
		//p := r.where[0]
		//for i := 1; i < len(r.where); i++ {
		//	p = p.And(r.where[i])
		//}
		//
		//if err := r.buildExpression(p); err != nil {
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

//func (r *Selector[T]) buildExpression(expr Expression) error {
//	switch exp := expr.(type) {
//	case nil:
//		return nil
//	case Predicate:
//		// 在这里构建p
//		// p.left构建好
//		// p.Op
//		// p.right
//		r.sb.WriteByte('(')
//		if err := r.buildExpression(exp.left); err != nil {
//			return err
//		}
//		//r.sb.WriteByte(' ')
//		if exp.op.String() != "" {
//			r.sb.WriteString(exp.op.String())
//			r.sb.WriteByte(' ')
//		}
//
//		if err := r.buildExpression(exp.right); err != nil {
//			return err
//		}
//		r.sb.WriteByte(')')
//	case Column:
//		//r.sb.WriteByte('`')
//		//fd, exist := r.model.FieldMap[exp.name]
//		//if !exist {
//		//	return errs.NewErrUnknownField(exp.name)
//		//}
//		//r.sb.WriteString(fd.ColName)
//		//r.sb.WriteByte('`')
//		//r.sb.WriteByte(' ')
//		exp.alias = ""
//		return r.buildColumn(exp)
//	case Op:
//	case value:
//		r.addArg(exp.val)
//		//r.args = append(r.args, raw.(value).val)
//		r.sb.WriteByte('?')
//	case RawExpr:
//		r.sb.WriteByte('(')
//		r.sb.WriteString(exp.raw)
//		r.addArg(exp.args...)
//		r.sb.WriteByte(')')
//	default:
//		return errs.NewErrUnsupportedExpression(expr)
//	}
//	return nil
//}

func (s *Selector[T]) buildColumns() (err error) {
	if len(s.cols) > 0 {
		//r.sb.WriteByte(' ')
		//for i, col := range r.cols {
		//	r.sb.WriteByte('`')
		//	r.sb.WriteString(col)
		//	r.sb.WriteByte('`')
		//	if i < len(r.cols)-1 {
		//		r.sb.WriteByte(',')
		//	}
		//	r.sb.WriteByte(' ')
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
			//r.sb.WriteByte(' ')
		}
		s.sb.WriteByte(' ')
	} else {
		s.sb.WriteString(" * ")
	}
	return nil
}

//func (r *Selector[T]) buildColumn(column Column) error {
//	r.sb.WriteByte('`')
//	fd, exist := r.model.FieldMap[column.name]
//	if !exist {
//		return errs.NewErrUnknownField(column.name)
//	}
//	r.sb.WriteString(fd.ColName)
//	r.sb.WriteByte('`')
//
//	if column.alias != "" {
//		r.sb.WriteByte(' ')
//		r.sb.WriteString("AS")
//		r.sb.WriteByte(' ')
//		r.sb.WriteByte('`')
//		r.sb.WriteString(column.alias)
//		r.sb.WriteByte('`')
//	}
//
//	//r.sb.WriteByte(' ')
//	return nil
//}
//
//func (r *Selector[T]) addArg(vals ...any) *Selector[T] {
//	if len(vals) == 0 {
//		return nil
//	}
//	if r.args == nil {
//		r.args = make([]any, 0, 4) // 给定预估容量，避免频繁扩容
//	}
//	r.args = append(r.args, vals...)
//	return r
//}

// select最简实现
//func (r *Selector[T]) Select(cols ...string) *Selector[T] {
//	r.cols = cols
//	return r
//}
//
//func (r *Selector[T]) Select(col string) *Selector[T] {
//	r.col = col
//	return r
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
//func (r Selector[T]) Get(ctx context.Context) (*T, error) {
//	q, err := r.Build()
//	if err != nil {
//		return nil, err
//	}
//
//	//var db *sql.DB
//	db := r.db.db
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
//	//r.model.FieldMap
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
//		fd, ok := r.model.ColumnMap[c]
//		if !ok {
//			return nil, errs.NewErrUnknownColumn(c)
//		}
//		val := reflect.New(fd.Typ)
//		vals = append(vals, val.Interface())
//		valElems = append(valElems, val.Elem())
//
//		//for _, fd := range r.model.FieldMap {
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
//	//	fd, ok := r.model.ColumnMap[c]
//	//	if !ok {
//	//		return nil, errs.NewErrUnknownColumn(c)
//	//	}
//	//	tpValue.Elem().FieldByName(fd.GoName).Set(valElems[i])
//	//	//tpValue.Elem().FieldByName(fd.GoName).Set(reflect.ValueOf(vals[i]).Elem())
//	//	//for _, fd := range r.model.FieldMap {
//	//	//	if fd.ColName == c {
//	//	//		tpValue.Elem().FieldByName(fd.GoName).Set(reflect.ValueOf(vals[i]).Elem())
//	//	//	}
//	//	//}
//	//}
//
//	tpValueElem := reflect.ValueOf(tp).Elem()
//	for i, c := range cs {
//		fd, ok := r.model.ColumnMap[c]
//		if !ok {
//			return nil, errs.NewErrUnknownColumn(c)
//		}
//		tpValueElem.FieldByName(fd.GoName).Set(valElems[i])
//		//tpValue.Elem().FieldByName(fd.GoName).Set(reflect.ValueOf(vals[i]).Elem())
//		//for _, fd := range r.model.FieldMap {
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
//func (r *Selector[T]) Get(ctx context.Context) (result *T, err error) {
//	if r.model == nil {
//		r.model, err = r.r.Get(new(T))
//		if err != nil {
//			return nil, err
//		}
//	}
//	root := r.getHandler
//	for i := len(r.mdls) - 1; i >= 0; i-- {
//		root = r.mdls[i](root)
//	}
//	res := root(ctx, &QueryContext{
//		Type:    "SELECT",
//		Builder: r,
//		Model:   r.model,
//	})
//	//var t *T
//	//if val, ok := res.Result.(*T);ok {
//	//	t = val
//	//}
//	//return t, res.Err
//	if res.Result != nil {
//		return res.Result.(*T), res.Err
//	}
//	return nil, res.Err
//}

func (s *Selector[T]) Get(ctx context.Context) (result *T, err error) {
	if s.model == nil {
		s.model, err = s.r.Get(new(T))
		if err != nil {
			return nil, err
		}
	}
	res := get[T](ctx, s.sess, s.core, &QueryContext{
		Type:    "SELECT",
		Builder: s,
		Model:   s.model,
	})
	if res.Result != nil {
		return res.Result.(*T), res.Err
	}
	return nil, res.Err
}

// var _ Handler = (&Selector[any]{}).getHandler

func getHandler[T any](ctx context.Context, sess Session, c core, qc *QueryContext) *QueryResult {
	q, err := qc.Builder.Build()
	if err != nil {
		//return nil, err
		return &QueryResult{Err: err}
	}

	// 发起查询，处理结果集
	row, err := sess.queryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return &QueryResult{Err: err}
	}

	if !row.Next() {
		// 里面是否返回error，返回error和sql包一致吗？和GetMulti保持一致
		return &QueryResult{Err: ErrNoRows}
	}

	tp := new(T)
	err = c.creator(c.model, tp).SetColumns(row)
	return &QueryResult{
		Result: tp,
		Err:    err,
	}
}

//func (r *Selector[T]) getHandler(ctx context.Context, qc *QueryContext) *QueryResult {
//	q, err := r.Build()
//	if err != nil {
//		//return nil, err
//		return &QueryResult{Err: err}
//	}
//
//	//var db *sql.DB
//	sess := r.sess
//	// 发起查询，处理结果集
//	row, err := sess.queryContext(ctx, q.SQL, q.Args...)
//	if err != nil {
//		return &QueryResult{Err: err}
//	}
//
//	if !row.Next() {
//		// 里面是否返回error，返回error和sql包一致吗？和GetMulti保持一致
//		return &QueryResult{Err: ErrNoRows}
//		//return nil, ErrNoRows
//	}
//
//	tp := new(T)
//	//var creator valuer.Creator
//	err = r.creator(r.model, tp).SetColumns(row)
//	//return tp, err
//	return &QueryResult{
//		Result: tp,
//		Err:    err,
//	}
//
//	////r.model.FieldMap
//	//
//	//// 问题： 类型、顺序要匹配
//	//
//	//// select出来哪些列
//	//cs, err := row.Columns()
//	//if err != nil {
//	//	return nil, err
//	//}
//	//
//	//tp := new(T)
//	//vals := make([]any, 0, len(cs))
//	//address := reflect.ValueOf(tp).UnsafePointer()
//	//for _, c := range cs {
//	//	fd, ok := r.model.ColumnMap[c]
//	//	if !ok {
//	//		return nil, errs.NewErrUnknownColumn(c)
//	//	}
//	//	// 起始地址+偏移量
//	//	fdAddress := unsafe.Pointer(uintptr(address) + fd.Offset)
//	//
//	//	// 反射在特定地址上，创建特定类型实例，原本类型的指针类型；case：fd.Typ=int, val是*int
//	//	val := reflect.NewAt(fd.Typ, fdAddress)
//	//	vals = append(vals, val.Interface())
//	//}
//	//err = row.Scan(vals...)
//	//if err != nil {
//	//	return nil, err
//	//}
//
//	//return tp, nil
//}

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
