package querylog

import (
	"awesomeProject1/orm"
	"context"
	"log"
)

type MiddlewareBuilder struct {
	logFunc func(query string, args []any)
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{logFunc: func(query string, args []any) {
		log.Println(query, args) // 未处理敏感数据
	}}
}

func (m *MiddlewareBuilder) LogFunc(fn func(query string, args []any)) *MiddlewareBuilder {
	m.logFunc = fn
	return m
}

func (m MiddlewareBuilder) Builder() orm.Middleware {
	return func(next orm.Handler) orm.Handler {
		return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			q, err := qc.Builder.Build()
			if err != nil {
				log.Println("构造sql出错", err)
				return &orm.QueryResult{Err: err}
			}
			//log.Println("sql: %s, args: %v\n", q.SQL, q.Args)
			m.logFunc(q.SQL, q.Args)
			result := next(ctx, qc)
			if result.Err != nil {
				m.logFunc(q.SQL, q.Args)
			}
			m.logFunc(q.SQL, q.Args)
			return result
		}
	}
}
