package orm

import (
	"context"
	"database/sql"
)

type RawQuerier[T any] struct {
	core
	Session
	sql  string
	args []any
}

func (r RawQuerier[T]) Build() (*Query, error) {
	return &Query{
		SQL:  r.sql,
		Args: r.args,
	}, nil
}

func RawQuery[T any](sess Session, query string, args ...any) *RawQuerier[T] {
	return &RawQuerier[T]{
		Session: sess,
		sql:     query,
		args:    args,
		core:    sess.getCore(),
	}
}

func (i RawQuerier[T]) Exec(ctx context.Context) Result {
	var err error
	i.model, err = i.r.Get(new(T))
	if err != nil {
		return Result{
			err: err,
		}
	}

	res := exec[T](ctx, i.Session, i.core, &QueryContext{
		Type:    "RAW",
		Builder: i,
		Model:   i.model,
	})
	var sqlRes sql.Result
	if res.Result != nil {
		sqlRes = res.Result.(sql.Result)
	}
	return Result{
		res: sqlRes, err: res.Err,
	}
}

func (r RawQuerier[T]) Get(ctx context.Context) (result *T, err error) {
	if r.model == nil {
		r.model, err = r.r.Get(new(T))
		if err != nil {
			return nil, err
		}
	}
	//var root Handler = func(ctx context.Context, qc *QueryContext) *QueryResult {
	//		return getHandler[T](ctx, r.Session, r.core, &QueryContext{
	//			Type:    "RAW",
	//			Builder: r,
	//			Model:   r.model,
	//		})
	//}
	//for i := len(r.mdls) - 1; i >= 0; i-- {
	//	root = r.mdls[i](root)
	//}
	//res := root(ctx, &QueryContext{
	//	Type:    "RAW",
	//	Builder: r,
	//	Model:   r.model,
	//})
	////var t *T
	////if val, ok := res.Result.(*T);ok {
	////	t = val
	////}
	////return t, res.Err
	res := get[T](ctx, r.Session, r.core, &QueryContext{
		Type:    "Raw",
		Builder: r,
		Model:   r.model,
	})
	if res.Result != nil {
		return res.Result.(*T), res.Err
	}
	return nil, res.Err
}

func get[T any](ctx context.Context, sess Session, c core, qc *QueryContext) *QueryResult {
	var err error
	if c.model == nil {
		c.model, err = c.r.Get(new(T))
		if err != nil {
			return &QueryResult{
				Err: err,
			}
		}
	}
	var root Handler = func(ctx context.Context, qc *QueryContext) *QueryResult {
		return getHandler[T](ctx, sess, c, &QueryContext{
			Type:    qc.Type,
			Builder: qc.Builder,
			Model:   c.model,
		})
	}
	for i := len(c.mdls) - 1; i >= 0; i-- {
		root = c.mdls[i](root)
	}
	return root(ctx, &QueryContext{
		Type:    qc.Type,
		Builder: qc.Builder,
		Model:   c.model,
	})
}

func (r RawQuerier[T]) GetMulti(ctx context.Context) ([]*T, error) {
	//TODO implement me
	panic("implement me")
}
