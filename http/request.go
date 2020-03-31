package http

import (
	"github.com/valyala/fasthttp"
)

type Requests interface {
	InternalServerErrorHandler(ctx *fasthttp.RequestCtx, err error)
	PanicHandler(ctx *fasthttp.RequestCtx, params interface{})
	MethodNotAllowed(ctx *fasthttp.RequestCtx)
	NotFound(ctx *fasthttp.RequestCtx)
	BasicRequest(ctx *fasthttp.RequestCtx)
	AdminRequest(ctx *fasthttp.RequestCtx)
}
