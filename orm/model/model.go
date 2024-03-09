package model

import (
	"awesomeProject1/orm/internal/errs"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

const (
	tagColumn = "column"
)

type RegistryInterface interface {
	Get(val any) (*Model, error)
	Registry(val any, opts ...Option) (*Model, error)
}

type TableName interface {
	TableName() string
}

type Model struct {
	TableName string
	// 字段名到字段的映射
	FieldMap map[string]*Field
	// 列名到字段的映射
	ColumnMap map[string]*Field
}

type Option func(m *Model) error

type Field struct {
	GoName  string
	ColName string
	Typ     reflect.Type // 字段类型

	// 字段相对于结构体本身的偏移量
	Offset uintptr
}

//var models = map[reflect.Type]*Model{}

//// 法一，并发工具
//type Registry struct {
//	lock   sync.RWMutex
//	models map[reflect.Type]*Model
//}
//
//func newRegistry() (*Registry, error) {
//	res := &Registry{
//		models: make(map[reflect.Type]*Model, 64),
//	}
//	return res, nil
//}

func MustNewRegistry() *Registry {
	res, err := newRegistry()
	if err != nil {
		panic(err)
	}
	return res
}

//func (r *Registry) Get1(val any) (m *Model, err error) {
//	Typ := reflect.TypeOf(val)
//	// double check
//	r.lock.RLock()
//	var ok bool
//	m, ok = r.models[Typ]
//	r.lock.RUnlock()
//	if ok {
//		return m, nil
//	}
//	r.lock.Lock()
//	defer r.lock.Unlock()
//
//	m, ok = r.models[Typ]
//	if ok {
//		return m, nil
//	}
//
//	m, err = r.Registry(val)
//	if err != nil {
//		return nil, err
//	}
//	r.models[Typ] = m
//
//	//if !ok {
//	//	m, err = r.Registry(val)
//	//	if err != nil {
//	//		return nil, err
//	//	}
//	//	r.models[Typ] = m
//	//}
//	return m, nil
//}
//
//var defaultRegistry = &Registry{
//	models: map[reflect.Type]*Model{},
//}

// 法二，
type Registry struct {
	models sync.Map
}

func newRegistry() (*Registry, error) {
	res := &Registry{}
	return res, nil
}

func (r *Registry) Get(val any) (*Model, error) {
	typ := reflect.TypeOf(val)
	m, ok := r.models.Load(typ)
	if ok {
		return m.(*Model), nil
	}
	//var err error
	m, err := r.Registry(val)
	if err != nil {
		return nil, err
	}
	//r.models.Store(Typ, m) // 同时解析会有覆盖的小问题
	return m.(*Model), err
}

// 限制只能使用一级指针
func (r *Registry) Registry(entity any, opts ...Option) (*Model, error) {
	typ := reflect.TypeOf(entity)
	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	elemTyp := typ.Elem()
	numField := elemTyp.NumField()
	fieldMap := make(map[string]*Field, numField)
	columnMap := make(map[string]*Field, numField)
	for i := 0; i < numField; i++ {
		fd := elemTyp.Field(i)
		pair, err := r.parseTag(fd.Tag)
		if err != nil {
			return nil, err
		}
		colName := pair[tagColumn]
		if colName == "" {
			// 用户未设置列名
			colName = underscoreName(fd.Name)
		}
		field := &Field{
			GoName:  fd.Name,
			ColName: colName,
			Typ:     fd.Type,
			Offset:  fd.Offset,
		}
		fieldMap[fd.Name] = field
		columnMap[colName] = field
	}

	var tableName string
	if tbl, ok := entity.(TableName); ok {
		tableName = tbl.TableName()
	}

	if tableName == "" {
		tableName = underscoreName(elemTyp.Name())
	}

	res := &Model{
		TableName: tableName,
		FieldMap:  fieldMap,
		ColumnMap: columnMap,
	}

	for _, opt := range opts {
		err := opt(res)
		if err != nil {
			return nil, err
		}
	}

	r.models.Store(typ, res)

	return res, nil
}

func WithTableName(tableName string) Option {
	return func(m *Model) error {
		m.TableName = tableName
		//if TableName == ""{
		//	return err
		//}
		return nil
	}
}

func WithColumnName(field string, colName string) Option {
	return func(m *Model) error {
		fd, ok := m.FieldMap[field]
		if !ok {
			return errs.NewErrUnknownField(field)
		}
		fd.ColName = colName
		return nil
	}
}

type User struct {
	ID uint64 `orm:"column=id"`
}

func (r *Registry) parseTag(tag reflect.StructTag) (map[string]string, error) {
	ormTag, ok := tag.Lookup("orm")
	if !ok {
		return map[string]string{}, nil
	}
	pairs := strings.Split(ormTag, ",")
	res := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		segs := strings.Split(pair, "=")
		if len(segs) != 2 {
			return nil, errs.NewErrInvalidTagContent(pair)
		}
		key := segs[0]
		res[key] = segs[1]
	}
	return res, nil
}

// underscoreName 驼峰转字符串命名
func underscoreName(tableName string) string {
	var buf []byte
	for i, v := range tableName {
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}

	}
	return string(buf)
}
