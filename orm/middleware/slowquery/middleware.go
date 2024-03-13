package slowquery

import (
	"awesomeProject1/orm"
	"context"
	"log"
	"time"
)

type MiddlewareBuilder struct {
	threshold time.Duration
	logFunc   func(query string, args []any)
}

// 100ms
func NewMiddlewareBuilder(threshold time.Duration) *MiddlewareBuilder {
	return &MiddlewareBuilder{logFunc: func(query string, args []any) {
		log.Println(query, args) // 未处理敏感数据
	},
		threshold: threshold,
	}
}

func (m *MiddlewareBuilder) Build() orm.Middleware {
	return func(next orm.Handler) orm.Handler {
		return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			startTime := time.Now()
			defer func() {
				duration := time.Since(startTime)
				if duration <= m.threshold {
					return
				}
				q, err := qc.Builder.Build()
				if err == nil {
					m.logFunc(q.SQL, q.Args)
				}
			}()
			return next(ctx, qc)
		}
	}
}
