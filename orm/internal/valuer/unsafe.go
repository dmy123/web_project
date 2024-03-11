package valuer

import (
	"awesomeProject1/orm/internal/errs"
	"awesomeProject1/orm/model"
	"database/sql"
	"reflect"
	"unsafe"
)

type unsafeValue struct {
	model *model.Model
	val   any // T的指针
	// 起始地址
	address unsafe.Pointer
}

func NewUnsafeValue(model *model.Model, val any) Value {
	return &unsafeValue{
		model:   model,
		val:     val,
		address: reflect.ValueOf(val).UnsafePointer(),
	}
}

var _ Creator = NewUnsafeValue

func (u unsafeValue) Field(name string) (any, error) {
	fd, ok := u.model.FieldMap[name]
	if !ok {
		return nil, errs.NewErrUnknownField(name)
	}
	fdAddress := unsafe.Pointer(uintptr(u.address) + fd.Offset)

	// 反射在特定地址上，创建特定类型实例，原本类型的指针类型；case：fd.Typ=int, val是*int
	val := reflect.NewAt(fd.Typ, fdAddress)
	return val.Elem().Interface(), nil
}

func (u unsafeValue) SetColumns(rows *sql.Rows) error {
	// select出来哪些列
	cs, err := rows.Columns()
	if err != nil {
		return err
	}

	vals := make([]any, 0, len(cs))
	//address := reflect.ValueOf(u.val).UnsafePointer()
	for _, c := range cs {
		fd, ok := u.model.ColumnMap[c]
		if !ok {
			return errs.NewErrUnknownColumn(c)
		}
		// 起始地址+偏移量
		fdAddress := unsafe.Pointer(uintptr(u.address) + fd.Offset)

		// 反射在特定地址上，创建特定类型实例，原本类型的指针类型；case：fd.Typ=int, val是*int
		val := reflect.NewAt(fd.Typ, fdAddress)
		vals = append(vals, val.Interface())
	}
	err = rows.Scan(vals...)
	if err != nil {
		return err
	}

	return nil
}
