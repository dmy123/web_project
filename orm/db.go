package orm

import (
	"awesomeProject1/orm/internal/valuer"
	"awesomeProject1/orm/model"
	"database/sql"
)

type DBOption func(db *DB)

type DB struct {
	r       *model.Registry
	db      *sql.DB
	creator valuer.Creator
}

func Open(driver string, dataSourceName string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open(driver, dataSourceName)
	if err != nil {
		return nil, err
	}
	return OpenDB(db, opts...)
	//res := &DB{
	//	r: MustNewRegistry(),
	//	db: db,
	//}
	//for _, opt := range opts {
	//	opt(res)
	//}
	//return res, nil
}

func MustOpen(driver string, dataSourceName string, opts ...DBOption) *DB {
	res, err := Open(driver, dataSourceName, opts...)
	if err != nil {
		panic(err)
	}
	return res
}

func OpenDB(db *sql.DB, opts ...DBOption) (*DB, error) {
	res := &DB{
		r:       model.MustNewRegistry(),
		db:      db,
		creator: valuer.NewUnsafeValue,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

func DBWithRegistry(r model.Registry) DBOption {
	return func(db *DB) {
		db.r = &r
	}
}

func DBUseReflect() DBOption {
	return func(db *DB) {
		db.creator = valuer.NewReflectValue
	}
}
