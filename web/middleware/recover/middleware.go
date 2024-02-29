package recover

import "awesomeProject1/web"

type MiddlewareBuilder struct {
	StatusCode int
	Data       []byte
	LogFunc    func(ctx *web.Context)
}

func (m MiddlewareBuilder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			defer func() {
				if err := recover(); err != nil {
					ctx.RespStatusCode = m.StatusCode
					ctx.RespData = m.Data
					m.LogFunc(ctx)
				}
			}()
			next(ctx)
		}
	}
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{}
}
