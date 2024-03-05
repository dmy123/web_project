package orm

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

type Registry interface {
	Get(val any) (*Model, error)
	Registry(val any, opts ...ModelOption) (*Model, error)
}

type Model struct {
	tableName string
	fields    map[string]*Field
}

type ModelOption func(m *Model) error

type Field struct {
	colName string
}

//var models = map[reflect.Type]*Model{}

//// 法一，并发工具
//type registry struct {
//	lock   sync.RWMutex
//	models map[reflect.Type]*Model
//}
//
//func newRegistry() (*registry, error) {
//	res := &registry{
//		models: make(map[reflect.Type]*Model, 64),
//	}
//	return res, nil
//}

func mustNewRegistry() *registry {
	res, err := newRegistry()
	if err != nil {
		panic(err)
	}
	return res
}

//func (r *registry) Get1(val any) (m *Model, err error) {
//	typ := reflect.TypeOf(val)
//	// double check
//	r.lock.RLock()
//	var ok bool
//	m, ok = r.models[typ]
//	r.lock.RUnlock()
//	if ok {
//		return m, nil
//	}
//	r.lock.Lock()
//	defer r.lock.Unlock()
//
//	m, ok = r.models[typ]
//	if ok {
//		return m, nil
//	}
//
//	m, err = r.Registry(val)
//	if err != nil {
//		return nil, err
//	}
//	r.models[typ] = m
//
//	//if !ok {
//	//	m, err = r.Registry(val)
//	//	if err != nil {
//	//		return nil, err
//	//	}
//	//	r.models[typ] = m
//	//}
//	return m, nil
//}
//
//var defaultRegistry = &registry{
//	models: map[reflect.Type]*Model{},
//}

// 法二，
type registry struct {
	models sync.Map
}

func newRegistry() (*registry, error) {
	res := &registry{}
	return res, nil
}

func (r *registry) Get(val any) (*Model, error) {
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
	//r.models.Store(typ, m) // 同时解析会有覆盖的小问题
	return m.(*Model), err
}

// 限制只能使用一级指针
func (r *registry) Registry(entity any, opts ...ModelOption) (*Model, error) {
	typ := reflect.TypeOf(entity)
	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	elemTyp := typ.Elem()
	numField := elemTyp.NumField()
	fieldMap := make(map[string]*Field, numField)
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
		fieldMap[fd.Name] = &Field{
			colName: colName,
		}
	}

	var tableName string
	if tbl, ok := entity.(TableName); ok {
		tableName = tbl.TableName()
	}

	if tableName == "" {
		tableName = underscoreName(elemTyp.Name())
	}

	res := &Model{
		tableName: tableName,
		fields:    fieldMap,
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

func ModelWithTableName(tableName string) ModelOption {
	return func(m *Model) error {
		m.tableName = tableName
		//if tableName == ""{
		//	return err
		//}
		return nil
	}
}

func ModelWithColumnName(field string, colName string) ModelOption {
	return func(m *Model) error {
		fd, ok := m.fields[field]
		if !ok {
			return errs.NewErrUnknownField(field)
		}
		fd.colName = colName
		return nil
	}
}

type User struct {
	ID uint64 `orm:"column=id"`
}

func (r *registry) parseTag(tag reflect.StructTag) (map[string]string, error) {
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
