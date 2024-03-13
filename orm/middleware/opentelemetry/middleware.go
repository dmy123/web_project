package opentelemetry

import (
	"awesomeProject1/orm"
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const instrumentationName = "awesomeProject1/orm/middleware/opentelemetry"

type MiddlewareBuilder struct {
	Tracer trace.Tracer
}

func (m MiddlewareBuilder) Build() orm.Middleware {
	if m.Tracer == nil {
		m.Tracer = otel.GetTracerProvider().Tracer(instrumentationName)
	}
	return func(next orm.Handler) orm.Handler {
		return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			tbl := qc.Model.TableName
			_, span := m.Tracer.Start(ctx, fmt.Sprintf("%s-%s", qc.Type, tbl))
			defer func() {
				span.End()
			}()

			q, _ := qc.Builder.Build()
			if q != nil {
				span.SetAttributes(attribute.String("sql", q.SQL))
			}

			span.SetAttributes(attribute.String("table", tbl))
			span.SetAttributes(attribute.String("component", "orm"))
			res := next(ctx, qc)
			if res.Err != nil {
				span.RecordError(res.Err)
			}
			return res
		}
	}
}
