package orm

import (
	"awesomeProject1/orm/internal/errs"
	"reflect"
	"sync"
	"unicode"
)

type model struct {
	tableName string
	fields    map[string]*field
}

type field struct {
	colName string
}

//var models = map[reflect.Type]*model{}

// 法一，并发工具
type registry struct {
	lock   sync.RWMutex
	models map[reflect.Type]*model
}

func newRegistry() (*registry, error) {
	res := &registry{
		models: make(map[reflect.Type]*model, 64),
	}
	return res, nil
}

func mustNewRegistry() *registry {
	res, err := newRegistry()
	if err != nil {
		panic(err)
	}
	return res
}

func (r *registry) get(val any) (m *model, err error) {
	typ := reflect.TypeOf(val)
	// double check
	r.lock.RLock()
	var ok bool
	m, ok = r.models[typ]
	r.lock.RUnlock()
	if ok {
		return m, nil
	}
	r.lock.Lock()
	defer r.lock.Unlock()

	m, ok = r.models[typ]
	if ok {
		return m, nil
	}

	m, err = r.parseModel(val)
	if err != nil {
		return nil, err
	}
	r.models[typ] = m

	//if !ok {
	//	m, err = r.parseModel(val)
	//	if err != nil {
	//		return nil, err
	//	}
	//	r.models[typ] = m
	//}
	return m, nil
}

var defaultRegistry = &registry{
	models: map[reflect.Type]*model{},
}

//// 法二，
//type registry struct {
//	models sync.Map
//}
//
//func newRegistry() (*registry, error) {
//	res := &registry{}
//	return res, nil
//}
//
//func (r *registry) get1(val any) (*model, error) {
//	typ := reflect.TypeOf(val)
//	m, ok := r.models.Load(typ)
//	if ok {
//		return m.(*model), nil
//	}
//	//var err error
//	m, err := r.parseModel(val)
//	if err != nil {
//		return nil, err
//	}
//	r.models.Store(typ, m) // 同时解析会有覆盖的小问题
//	return m.(*model), err
//}

// 限制只能使用一级指针
func (r *registry) parseModel(entity any) (*model, error) {
	typ := reflect.TypeOf(entity)
	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	typ = typ.Elem()
	numField := typ.NumField()
	fieldMap := make(map[string]*field, numField)
	for i := 0; i < numField; i++ {
		fd := typ.Field(i)
		fieldMap[fd.Name] = &field{
			colName: underscoreName(fd.Name),
		}
	}
	return &model{
		tableName: underscoreName(typ.Name()),
		fields:    fieldMap,
	}, nil
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
