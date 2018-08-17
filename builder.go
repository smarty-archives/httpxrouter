package httpxrouter

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/smartystreets/httpx/middleware"
)

func New(options ...Option) http.Handler {
	builder := new(routesBuilder)
	builder.router = &httprouter.Router{
		RedirectTrailingSlash:  false,
		RedirectFixedPath:      false,
		HandleMethodNotAllowed: true,
		HandleOPTIONS:          true,
	}
	builder.apply(options)
	return builder.build()
}

type routesBuilder struct {
	handlers []NestingHandler
	router   *httprouter.Router
}

func (builder *routesBuilder) apply(options []Option) {
	for _, option := range options {
		option(builder)
	}
}

func (this *routesBuilder) build() http.Handler {
	return chainN(append(this.handlers, middleware.NewNestableHandler(this.router)))
}

func routeFunc(handler http.Handler) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(response http.ResponseWriter, request *http.Request, _ httprouter.Params) {
		handler.ServeHTTP(response, request)
	}
}

func chainN(handlers []NestingHandler) NestingHandler {
	for x := 0; x < len(handlers)-1; x++ {
		handlers[x].Install(handlers[x+1])
	}
	return handlers[0]
}
