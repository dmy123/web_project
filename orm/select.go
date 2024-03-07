package orm

import (
	"awesomeProject1/orm/internal/errs"
	"awesomeProject1/orm/internal/valuer"
	"context"
	"strings"
)

type Selector[T any] struct {
	table string
	model *Model
	where []Predicate
	sb    *strings.Builder
	args  []any
	db    *DB
	//r *registry
}

//func (db *DB) NewSelector[T any]() *Selector[T] {
//	return &Selector[T]{
//		sb: &strings.Builder{},
//		db: db,
//	}
//}

func NewSelector[T any](db *DB) *Selector[T] {
	return &Selector[T]{
		sb: &strings.Builder{},
		db: db,
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
	s.model, err = s.db.r.Registry(new(T))
	if err != nil {
		return nil, err
	}
	s.sb.WriteString("SELECT * FROM ")
	// 反射拿到表名
	if s.table == "" {
		//var t T
		//typ := reflect.TypeOf(t)
		//s.table = typ.Name()
		s.sb.WriteByte('`')
		//s.sb.WriteString(typ.Name())
		s.sb.WriteString(s.model.tableName)
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
	switch exp := expr.(type) {
	case nil:
		return nil
	case Predicate:
		// 在这里构建p
		// p.left构建好
		// p.Op
		// p.right
		s.sb.WriteByte('(')
		if err := s.buildExpression(exp.left); err != nil {
			return err
		}
		//s.sb.WriteByte(' ')
		s.sb.WriteString(exp.op.String())
		s.sb.WriteByte(' ')

		if err := s.buildExpression(exp.right); err != nil {
			return err
		}
		s.sb.WriteByte(')')
	case Column:
		s.sb.WriteByte('`')
		fd, exist := s.model.fieldMap[exp.name]
		if !exist {
			return errs.NewErrUnknownField(exp.name)
		}
		s.sb.WriteString(fd.colName)
		s.sb.WriteByte('`')
		s.sb.WriteByte(' ')
	case Op:
	case value:
		s.addArg(exp.val)
		//s.args = append(s.args, expr.(value).val)
		s.sb.WriteByte('?')
	default:
		return errs.NewErrUnsupportedExpression(expr)
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
//	//s.model.fieldMap
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
//		fd, ok := s.model.columnMap[c]
//		if !ok {
//			return nil, errs.NewErrUnknownColumn(c)
//		}
//		val := reflect.New(fd.typ)
//		vals = append(vals, val.Interface())
//		valElems = append(valElems, val.Elem())
//
//		//for _, fd := range s.model.fieldMap {
//		//	if fd.colName == c {
//		//		// 反射创建新的实例
//		//		val := reflect.New(fd.typ)
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
//	//	fd, ok := s.model.columnMap[c]
//	//	if !ok {
//	//		return nil, errs.NewErrUnknownColumn(c)
//	//	}
//	//	tpValue.Elem().FieldByName(fd.goName).Set(valElems[i])
//	//	//tpValue.Elem().FieldByName(fd.goName).Set(reflect.ValueOf(vals[i]).Elem())
//	//	//for _, fd := range s.model.fieldMap {
//	//	//	if fd.colName == c {
//	//	//		tpValue.Elem().FieldByName(fd.goName).Set(reflect.ValueOf(vals[i]).Elem())
//	//	//	}
//	//	//}
//	//}
//
//	tpValueElem := reflect.ValueOf(tp).Elem()
//	for i, c := range cs {
//		fd, ok := s.model.columnMap[c]
//		if !ok {
//			return nil, errs.NewErrUnknownColumn(c)
//		}
//		tpValueElem.FieldByName(fd.goName).Set(valElems[i])
//		//tpValue.Elem().FieldByName(fd.goName).Set(reflect.ValueOf(vals[i]).Elem())
//		//for _, fd := range s.model.fieldMap {
//		//	if fd.colName == c {
//		//		tpValue.Elem().FieldByName(fd.goName).Set(reflect.ValueOf(vals[i]).Elem())
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
	db := s.db.db
	// 发起查询，处理结果集
	row, err := db.QueryContext(ctx, q.SQL, q.Args...)
	if err != nil {
		return nil, err
	}

	if !row.Next() {
		// 里面是否返回error，返回error和sql包一致吗？和GetMulti保持一致
		return nil, ErrNoRows
	}

	tp := new(T)
	var creator valuer.Creator
	err = creator(tp).SetColumns(row)
	return tp, err

	////s.model.fieldMap
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
	//	fd, ok := s.model.columnMap[c]
	//	if !ok {
	//		return nil, errs.NewErrUnknownColumn(c)
	//	}
	//	// 起始地址+偏移量
	//	fdAddress := unsafe.Pointer(uintptr(address) + fd.offset)
	//
	//	// 反射在特定地址上，创建特定类型实例，原本类型的指针类型；case：fd.typ=int, val是*int
	//	val := reflect.NewAt(fd.typ, fdAddress)
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
	db := s.db.db
	// 发起查询，处理结果集
	rows, err := db.QueryContext(ctx, q.SQL, q.Args)

	for rows.Next() {

	}
	return nil, nil
}
