package safedml

import (
	"awesomeProject1/orm"
	"context"
	"errors"
	"strings"
)

// 强制查询语句
// 1、SELECT UPDATE DELETE必须带WHERE
// 2、UPDATE 和 DELETE 必须带 WHERE
type MiddlewareBuilder struct {
}

func NewMiddlewareBuilder() {

}

func (m *MiddlewareBuilder) Build() orm.Middleware {
	return func(next orm.Handler) orm.Handler {
		return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			if qc.Type == "SELECT" || qc.Type == "INSERT" {
				return next(ctx, qc)
			}
			q, err := qc.Builder.Build()
			if err != nil {
				return &orm.QueryResult{
					Err: err,
				}
			}
			if strings.Contains(q.SQL, "WHERE") {
				return &orm.QueryResult{
					Err: errors.New("不准执行没有 WHERE 的 DELETE or UPDATE 语句"),
				}
			}
			return next(ctx, qc)
		}
	}
}
