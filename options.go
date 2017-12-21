package httpxrouter

import (
	"net/http"
	"strings"
)

type Option func(*routesBuilder)

func NotFound(handler http.Handler) Option {
	return func(builder *routesBuilder) {
		builder.router.NotFound = handler
	}
}

func MethodNotAllowed(handler http.Handler) Option {
	return func(builder *routesBuilder) {
		builder.router.MethodNotAllowed = handler
	}
}

func Panic(callback func(http.ResponseWriter, *http.Request, interface{})) Option {
	return func(builder *routesBuilder) {
		builder.router.PanicHandler = callback
	}
}

func Prepend(handlers ...nestingHandler) Option {
	return func(builder *routesBuilder) {
		builder.handlers = append(handlers, builder.handlers...)
	}
}

func HEAD(path string, handlers ...nestingHandler) Option {
	return Register("HEAD", path, handlers...)
}

func OPTIONS(path string, handlers ...nestingHandler) Option {
	return Register("OPTIONS", path, handlers...)
}

func GET(path string, handlers ...nestingHandler) Option {
	return Register("GET", path, handlers...)
}

func PUT(path string, handlers ...nestingHandler) Option {
	return Register("PUT", path, handlers...)
}

func POST(path string, handlers ...nestingHandler) Option {
	return Register("POST", path, handlers...)
}

func DELETE(path string, handlers ...nestingHandler) Option {
	return Register("DELETE", path, handlers...)
}

func Register(methods, paths string, handlers ...nestingHandler) Option {
	return func(builder *routesBuilder) {
		handler := routeFunc(chainN(handlers))
		for _, method := range strings.Split(methods, "|") {
			for _, path := range strings.Split(paths, "|") {
				builder.router.Handle(method, path, handler)
			}
		}
	}
}
