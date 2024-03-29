package orm

import (
	"context"
)

// Querier[T any] 用于select语句
type Querier[T any] interface {
	Get(ctx context.Context) (*T, error)
	GetMulti(ctx context.Context) ([]*T, error)
}

// Executor 用于insert、delete、update
type Executor interface {
	Exec(ctx context.Context) Result
}

type QueryBuilder interface {
	Build() (*Query, error)
}

type Query struct {
	SQL  string
	Args []any
}
