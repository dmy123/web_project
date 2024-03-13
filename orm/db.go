package orm

import (
	"awesomeProject1/orm/internal/errs"
	"awesomeProject1/orm/internal/valuer"
	"awesomeProject1/orm/model"
	"context"
	"database/sql"
	"database/sql/driver"
	"log"
	"time"
)

type DBOption func(db *DB)

type DB struct {
	core
	//r       *model.Registry
	db *sql.DB

	//creator valuer.Creator
	//dialect Dialect
}

// Wait 会等待数据库连接
// 注意只能用于测试
func (db *DB) Wait() error {
	err := db.db.Ping()
	for err == driver.ErrBadConn {
		log.Printf("等待数据库启动...")
		err = db.db.Ping()
		time.Sleep(time.Second)
	}
	return err
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

type txKey struct {
}

//// BeginTxV2 事务已被提交情况
//func (db *DB) BeginTxV2(ctx context.Context, opts *sql.TxOptions) (context.Context, *Tx, error) {
//	val := ctx.Value(txKey{})
//	tx, ok := val.(*Tx)
//	if ok {
//		return ctx, tx, nil
//	}
//	tx, err := db.BeginTx(ctx, opts)
//	if err != nil {
//		return ctx, nil, err
//	}
//	ctx = context.WithValue(ctx, txKey{}, tx)
//	return ctx, tx, nil
//}
//
//// BeginTxV3 必须上层链路开事务
//func (db *DB) BeginTxV3(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
//	val := ctx.Value(txKey{})
//	tx, ok := val.(*Tx)
//	if ok {
//		return tx, nil
//	}
//	return nil, errors.New("未开事务")
//}

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

func DBWithMiddleware(mdls ...Middleware) DBOption {
	return func(db *DB) {
		db.mdls = append(db.mdls, mdls...)
	}
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
