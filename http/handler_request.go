package http

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/DenisGubenko/ideasymbols/db"
	"github.com/DenisGubenko/ideasymbols/utils"
	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

const (
	headerAcceptApplicationJSONValue = `application/json`
	unknownStatusCode                = `Unknown Status Code`
)

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type answerResponse struct {
	Result map[string]interface{} `json:"result"`
	Error  errorResponse          `json:"error"`
}

type handlerRequest struct {
	storage db.Storage
}

func newHandlerRequest(storage db.Storage) Requests {
	return &handlerRequest{
		storage: storage,
	}
}

func (h *handlerRequest) sendJSONResponse(
	ctx *fasthttp.RequestCtx, input map[string]interface{},
	errorCode int, errorMsg *string) {
	h.defaultJSONCtx(ctx)
	if errorMsg == nil || len(*errorMsg) == 0 {
		errorMsgLocal := fasthttp.StatusMessage(errorCode)
		if errorMsgLocal == unknownStatusCode {
			errorCode = fasthttp.StatusInternalServerError
			errorMsgLocal = fasthttp.StatusMessage(errorCode)
		}
		errorMsg = &errorMsgLocal
	}
	ctx.SetStatusCode(errorCode)

	result := answerResponse{}
	if input != nil {
		result.Result = input
	}
	result.Error = errorResponse{
		Code:    errorCode,
		Message: *errorMsg,
	}
	answer, err := json.Marshal(result)
	if err != nil {
		h.shutdownErrorAnswer(ctx, err)
		return
	}
	ctx.SetBodyString(string(answer))
}

func (h *handlerRequest) InternalServerErrorHandler(ctx *fasthttp.RequestCtx, err error) {
	h.defaultJSONCtx(ctx)
	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	logrus.Errorf(
		fmt.Sprintf(errors.WithStack(err).Error()+
			`; route '%s'; request type '%s'; request header: '%s'; request body: '%s'`,
			string(ctx.Request.RequestURI()),
			string(ctx.Method()),
			string(ctx.Request.Header.Header()),
			string(ctx.Request.Body())))
	errorMsg := err.Error()
	result := answerResponse{
		Error: errorResponse{
			Code:    fasthttp.StatusInternalServerError,
			Message: errorMsg,
		},
	}
	answer, err := json.Marshal(result)
	if err != nil {
		h.shutdownErrorAnswer(ctx, err)
		return
	}
	ctx.SetBodyString(string(answer))
}

func (h *handlerRequest) PanicHandler(ctx *fasthttp.RequestCtx, params interface{}) {
	h.InternalServerErrorHandler(ctx, errors.New(fmt.Sprintf("%v", params)))
}

func (h *handlerRequest) MethodNotAllowed(ctx *fasthttp.RequestCtx) {
	if ok := h.echoAnswerIfHeaderAcceptNotEqualApplicationJSON(ctx, fasthttp.StatusMethodNotAllowed, nil); ok {
		return
	}
	h.defaultJSONCtx(ctx)
	code := fasthttp.StatusMethodNotAllowed
	h.sendJSONResponse(ctx, nil, code, nil)
}

func (h *handlerRequest) NotFound(ctx *fasthttp.RequestCtx) {
	if ok := h.echoAnswerIfHeaderAcceptNotEqualApplicationJSON(ctx, fasthttp.StatusNotFound, nil); ok {
		return
	}
	h.defaultJSONCtx(ctx)
	code := fasthttp.StatusNotFound
	h.sendJSONResponse(ctx, nil, code, nil)
}

func (h *handlerRequest) BasicRequest(ctx *fasthttp.RequestCtx) {
	result, err := h.storage.GetRandomOrderContent()
	if err != nil {
		ctx.SetBodyString(fmt.Sprintf(`%+v`, errors.WithStack(err)))
	}

	ctx.SetBodyString(result.Content)
}

func (h *handlerRequest) AdminRequest(ctx *fasthttp.RequestCtx) {
	count, orders, err := h.storage.GetStatisticsOrder()
	if err != nil {
		ctx.SetBodyString(fmt.Sprintf(`%+v`, errors.WithStack(err)))
	}

	content := `<html><body>`
	content += fmt.Sprintf(`Total: %d`, *count)

	for _, value := range orders {
		content += fmt.Sprintf(`<br>%s - %d`, value.Content, value.Counter)
	}

	content += `</body></html>`

	ctx.SetContentType("text/html")
	ctx.SetBody([]byte(content))
}

func (h *handlerRequest) echoAnswerIfHeaderAcceptNotEqualApplicationJSON(
	ctx *fasthttp.RequestCtx, code int, errMsg *string) bool {
	if !strings.Contains(string(ctx.Request.Header.Peek(`Accept`)), headerAcceptApplicationJSONValue) {
		h.defaultHTMLCtx(ctx)
		ctx.Response.SetStatusCode(code)
		if errMsg == nil {
			errorMsg := fasthttp.StatusMessage(code)
			errMsg = &errorMsg
		}
		ctx.Response.SetBodyString(*errMsg)
		return true
	}
	return false
}

func (h *handlerRequest) shutdownErrorAnswer(ctx *fasthttp.RequestCtx, err error) {
	msg, _ := json.Marshal(err.Error())
	ctx.SetBodyString(
		fmt.Sprintf(`{"error":{"code":%d,"message":"%s"}}`, fasthttp.StatusInternalServerError, msg))
}

func (h *handlerRequest) defaultHTMLCtx(ctx *fasthttp.RequestCtx) {
	ctx.Response.Reset()
	ctx.Response.Header.SetServer(utils.GetEnvVariable(`HTTP_SERVER_NAME`))
	ctx.SetContentType("text/html; charset=utf-8")
}

func (h *handlerRequest) defaultJSONCtx(ctx *fasthttp.RequestCtx) {
	ctx.Response.Reset()
	ctx.Response.Header.SetServer(utils.GetEnvVariable(`HTTP_SERVER_NAME`))
	ctx.SetContentType("application/json; charset=utf-8")
}
