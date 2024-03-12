package orm

import "context"

type QueryContext struct {
	Type    string // 标记增删改查
	Builder QueryBuilder
}

type QueryResult struct {
	// Result 在不同查询下类型不同，select可为*T或[]*T，其它是类型Result
	Result any
	Err    error
}

type Handler func(ctx context.Context, qc *QueryContext) *QueryResult

type Middleware func(next Handler) Handler
