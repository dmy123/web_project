package sql

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JsonColumn[T any] struct {
	Val   T
	Valid bool // 是否非nil
}

func (j *JsonColumn[T]) Scan(src any) error {
	//    int64
	//    float64
	//    bool
	//    []byte
	//    string
	//    time.Time
	//    nil - for NULL values
	var bs []byte
	switch data := src.(type) {
	case string:
		bs = []byte(data) // 可以考虑额外处理空字符串
	case []byte:
		bs = data // 可以考虑额外处理[]byte{}
	case nil:
		// 说明数据库存的就是nil
		return nil
	default:
		return errors.New("不支持类型")
	}
	err := json.Unmarshal(bs, &j.Val)
	if err == nil {
		j.Valid = true
	}
	return err
}

func (j JsonColumn[T]) Value() (driver.Value, error) {
	if !j.Valid {
		return nil, nil // nil 也是合法值
	}
	return json.Marshal(j.Val)
}
