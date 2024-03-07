package unsafe

import (
	"errors"
	"reflect"
	"unsafe"
)

type UnsafeAccessor struct {
	fields  map[string]FieldMeta
	address unsafe.Pointer
}

func NewUnsafeAccessor(entity any) *UnsafeAccessor {
	typ := reflect.TypeOf(entity)
	typ = typ.Elem()
	numFields := typ.NumField()
	fields := make(map[string]FieldMeta, numFields)
	for i := 0; i < numFields; i++ {
		fd := typ.Field(i)
		fields[fd.Name] = FieldMeta{Offset: fd.Offset, typ: fd.Type}
	}

	val := reflect.ValueOf(entity)
	return &UnsafeAccessor{
		fields:  fields,
		address: val.UnsafePointer(), // UnsafeAddress不稳定
	}
}

func (a *UnsafeAccessor) Field(field string) (any, error) {
	fd, ok := a.fields[field]
	if !ok {
		return nil, errors.New("非法字段")
	}
	fdAddress := unsafe.Pointer(uintptr(a.address) + fd.Offset)
	// 已知确切类型
	//return *(*int)(fdAddress), nil
	// 未知确切类型
	return reflect.NewAt(fd.typ, fdAddress).Elem().Interface(), nil
}

func (a *UnsafeAccessor) SetField(field string, val any) error {
	// 起始地址+字段量偏移
	fd, ok := a.fields[field]
	if !ok {
		return errors.New("非法字段")
	}
	fdAddress := unsafe.Pointer(uintptr(a.address) + fd.Offset)

	// 已知确切类型
	//*(*int)(fdAddress) = val.(int)
	//return nil
	reflect.NewAt(fd.typ, fdAddress).Elem().Set(reflect.ValueOf(val))
	return nil
}

type FieldMeta struct {
	Offset uintptr
	typ    reflect.Type
}
