package orm

import (
	"awesomeProject1/orm/internal/valuer"
	"awesomeProject1/orm/model"
	"context"
)

type core struct {
	model *model.Model

	dialect Dialect
	creator valuer.Creator
	r       *model.Registry

	mdls []Middleware
}

func execHandler(ctx context.Context, sess Session, c core, qc *QueryContext) *QueryResult {
	q, err := qc.Builder.Build()
	if err != nil {
		return &QueryResult{
			Result: Result{
				err: err,
			},
			Err: err,
		}
	}
	r, err := sess.execContext(ctx, q.SQL, q.Args...)
	return &QueryResult{
		Result: Result{
			err: err,
			res: r,
		},
		Err: err,
	}
}
