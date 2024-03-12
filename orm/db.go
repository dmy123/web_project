package orm

import (
	"awesomeProject1/orm/internal/errs"
	"awesomeProject1/orm/internal/valuer"
	"awesomeProject1/orm/model"
	"context"
	"database/sql"
)

type DBOption func(db *DB)

type DB struct {
	core
	//r       *model.Registry
	db *sql.DB

	//creator valuer.Creator
	//dialect Dialect
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
		core: core{
			r:       model.MustNewRegistry(),
			creator: valuer.NewUnsafeValue,
			dialect: DialectMySQL,
		},
		db: db,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res, nil
}

func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Tx{
		tx: tx,
		//db: db,
	}, nil
}

func (db *DB) getCore() core {
	return db.core
}

func (db *DB) DoTx(ctx context.Context,
	fn func(ctx context.Context, tx *Tx) error, opts *sql.TxOptions) (err error) {
	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return err
	}
	paniced := true

	defer func() {
		if paniced || err != nil {
			e := tx.Rollback()
			err = errs.NewErrFailedToRollbackTx(err, e, paniced)
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(ctx, tx)
	paniced = false

	return err
}

func (db *DB) queryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return db.db.QueryContext(ctx, query, args...)
}

func (db *DB) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return db.db.ExecContext(ctx, query, args...)
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

func DBWithDialect(d Dialect) DBOption {
	return func(db *DB) {
		db.dialect = d
	}
}
