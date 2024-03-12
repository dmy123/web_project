package web

// Middleware 函数式的责任链模式/函数式的洋葱模式
type Middleware func(next HandleFunc) HandleFunc

//type MiddlewareV1 interface {
//	Invoke(next HandleFunc) HandleFunc
//}
//
//// 拦截器设计
//type Interceptor interface {
//	Before(ctx *Context)
//	After(ctx *Context)
//	Surround(ctx *Context)
//}

// 集中式
//type Chain []HandleFunc
//
//type HandleFuncV1 func(ctx *Context) (next bool)
//
//
//type ChainV1 struct {
//	handlers []HandleFuncV1
//}
//
//func (c ChainV1) Run(ctx *Context)  {
//	for _, h := range c.handlers{
//		next := h(ctx)
//		if !next{
//			return
//		}
//	}
//}

//type Net struct {
//	handlers []HandleFuncV2
//}
//
//func (n Net) Run(ctx *Context) {
//	var wg sync.WaitGroup
//	for _, hdl := range n.handlers {
//		h := hdl
//		if h.concurrent{
//			wg.Add(1)
//			go func() {
//				h.Run(ctx)
//				wg.Done()
//			}()
//		}else {
//			h.Run(ctx)
//		}
//	}
//	wg.Wait()
//}
//
//type HandleFuncV2 struct {
//	concurrent bool
//	handlers []*HandleFuncV2
//}
//
//func (h HandleFuncV2) Run(ctx *Context)  {
//	for _, hdl := range n.handlers {
//		h := hdl
//		if h.concurrent{
//			wg.Add(1)
//			go func() {
//				h.Run(ctx)
//				wg.Done()
//			}()
//		}
//	}
//}
