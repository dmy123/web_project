package orm

import "context"

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

func RawQuery[T any](query string, args ...any) *RawQuerier[T] {
	return &RawQuerier[T]{
		sql:  query,
		args: args,
	}
}

func (r RawQuerier[T]) Exec(ctx context.Context) Result {
	//TODO implement me
	panic("implement me")
}

func (r RawQuerier[T]) Get(ctx context.Context) (result *T, err error) {
	//if s.model == nil {
	//	s.model, err = s.r.Get(new(T))
	//	if err != nil {
	//		return nil, err
	//	}
	//}
	//var root Handler = func(ctx context.Context, qc *QueryContext) *QueryResult {
	//		return getHandler[T](ctx, s.Session, s.core, &QueryContext{
	//			Type:    "RAW",
	//			Builder: s,
	//			Model:   s.model,
	//		})
	//}
	//for i := len(s.mdls) - 1; i >= 0; i-- {
	//	root = s.mdls[i](root)
	//}
	//res := root(ctx, &QueryContext{
	//	Type:    "RAW",
	//	Builder: s,
	//	Model:   s.model,
	//})
	////var t *T
	////if val, ok := res.Result.(*T);ok {
	////	t = val
	////}
	////return t, res.Err
	res := get[T](ctx, r.Session, r.core, r)
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
