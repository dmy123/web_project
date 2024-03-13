package nodelete

import (
	"awesomeProject1/orm"
	"context"
	"errors"
)

type MiddlewareBuilder struct {
}

func NewMiddlewareBuilder() {

}

func (m *MiddlewareBuilder) Build() orm.Middleware {
	return func(next orm.Handler) orm.Handler {
		return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			if qc.Type == "DELETE" {
				return &orm.QueryResult{
					Err: errors.New("禁止 Delete 语句"),
				}
			}
			return next(ctx, qc)
		}
	}
}
