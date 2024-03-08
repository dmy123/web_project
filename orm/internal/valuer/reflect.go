package valuer

import (
	"awesomeProject1/orm/internal/errs"
	"awesomeProject1/orm/model"
	"database/sql"

	"reflect"
)

type reflectValue struct {
	model *model.Model
	val   any // T的指针
}

func NewReflectValue(model *model.Model, val any) Value {
	return &reflectValue{
		model: model,
		val:   val,
	}
}

var _ Creator = NewReflectValue

func (r reflectValue) SetColumns(rows *sql.Rows) (err error) {
	cs, err := rows.Columns()
	if err != nil {
		return err
	}
	vals := make([]any, 0, len(cs))
	valElems := make([]reflect.Value, 0, len(cs))
	for _, c := range cs {
		fd, ok := r.model.ColumnMap[c]
		if !ok {
			return errs.NewErrUnknownColumn(c)
		}
		val := reflect.New(fd.Typ)
		vals = append(vals, val.Interface())
		valElems = append(valElems, val.Elem())

		//for _, fd := range s.model.fieldMap {
		//	if fd.colName == c {
		//		// 反射创建新的实例
		//		val := reflect.New(fd.typ)
		//		vals = append(vals, val.Interface())
		//	}
		//}
	}
	err = rows.Scan(vals...)
	if err != nil {
		return err
	}

	//tpValue := reflect.ValueOf(tp)
	//for i, c := range cs {
	//	fd, ok := s.model.columnMap[c]
	//	if !ok {
	//		return nil, errs.NewErrUnknownColumn(c)
	//	}
	//	tpValue.Elem().FieldByName(fd.goName).Set(valElems[i])
	//	//tpValue.Elem().FieldByName(fd.goName).Set(reflect.ValueOf(vals[i]).Elem())
	//	//for _, fd := range s.model.fieldMap {
	//	//	if fd.colName == c {
	//	//		tpValue.Elem().FieldByName(fd.goName).Set(reflect.ValueOf(vals[i]).Elem())
	//	//	}
	//	//}
	//}

	tpValueElem := reflect.ValueOf(r.val).Elem()
	for i, c := range cs {
		fd, ok := r.model.ColumnMap[c]
		if !ok {
			return errs.NewErrUnknownColumn(c)
		}
		tpValueElem.FieldByName(fd.GoName).Set(valElems[i])
		//tpValue.Elem().FieldByName(fd.goName).Set(reflect.ValueOf(vals[i]).Elem())
		//for _, fd := range s.model.fieldMap {
		//	if fd.colName == c {
		//		tpValue.Elem().FieldByName(fd.goName).Set(reflect.ValueOf(vals[i]).Elem())
		//	}
		//}
	}

	return nil
}
