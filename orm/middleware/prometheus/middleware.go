package prometheus

import (
	"awesomeProject1/orm"
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

type MiddlewareBuilder struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
}

func (m MiddlewareBuilder) Build() orm.Middleware {
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:      m.Name,
		Subsystem: m.Subsystem,
		Namespace: m.Namespace,
		Help:      m.Help,
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.90:  0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	}, []string{"type", "table"})
	prometheus.MustRegister(vector)
	return func(next orm.Handler) orm.Handler {
		return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			startTime := time.Now()
			defer func() {
				endTime := time.Now()
				typ := "unknown"
				// 原生查询才会走到这里
				tblName := "unknown"
				if qc.Model != nil {
					typ = qc.Model.TableName
					tblName = qc.Model.TableName
				}
				vector.WithLabelValues(typ, tblName).
					Observe(float64(endTime.Sub(startTime).Milliseconds()))
			}()
			return next(ctx, qc)
		}
	}
}
