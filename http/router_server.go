package http

import (
	"fmt"

	"github.com/DenisGubenko/ideasymbols/db"
	"github.com/DenisGubenko/ideasymbols/utils"
	"github.com/pkg/errors"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

type routerServer struct {
	storage  db.Storage
	requests Requests
	router   *fasthttprouter.Router
	server   *fasthttp.Server
}

func NewRouterServer(storage db.Storage) Server {
	return &routerServer{
		storage: storage,
	}
}

func (r *routerServer) Start() error {
	r.requests = newHandlerRequest(r.storage)

	r.initRoutes()

	if err := r.initHTTPServer(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *routerServer) Stop() (err error) {
	err = r.server.Shutdown()
	if err != nil {
		return err
	}

	r.server = nil
	r.router = nil
	r.requests = nil
	return
}

func (r *routerServer) initHTTPServer() (err error) {
	r.server = &fasthttp.Server{
		Handler: r.router.Handler,
	}

	err = r.server.ListenAndServe(fmt.Sprintf(`:%d`, utils.GetIntEnvVariable(`HTTP_PORT`)))
	if err != nil {
		return err
	}
	return nil
}

func (r *routerServer) initRoutes() {
	r.router = &fasthttprouter.Router{
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
		HandleOPTIONS:          true,
		NotFound:               r.requests.NotFound,
		MethodNotAllowed:       r.requests.MethodNotAllowed,
		PanicHandler:           r.requests.PanicHandler,
	}

	r.customRoutes()
}

func (r *routerServer) customRoutes() {
	r.router.GET("/request", r.requests.BasicRequest)
	r.router.GET("/admin/requests", r.requests.AdminRequest)
}
