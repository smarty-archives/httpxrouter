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

func Prepend(handlers ...NestingHandler) Option {
	return func(builder *routesBuilder) {
		builder.handlers = append(handlers, builder.handlers...)
	}
}

func HEAD(path string, handlers ...NestingHandler) Option {
	return Register("HEAD", path, handlers...)
}

func OPTIONS(path string, handlers ...NestingHandler) Option {
	return Register("OPTIONS", path, handlers...)
}

func GET(path string, handlers ...NestingHandler) Option {
	return Register("GET", path, handlers...)
}

func PUT(path string, handlers ...NestingHandler) Option {
	return Register("PUT", path, handlers...)
}

func POST(path string, handlers ...NestingHandler) Option {
	return Register("POST", path, handlers...)
}

func DELETE(path string, handlers ...NestingHandler) Option {
	return Register("DELETE", path, handlers...)
}

func Register(methods, paths string, handlers ...NestingHandler) Option {
	return func(builder *routesBuilder) {
		handler := chainN(handlers)
		for _, method := range strings.Split(methods, "|") {
			for _, path := range strings.Split(paths, "|") {
				builder.router.Handler(method, path, handler)
			}
		}
	}
}

func Compound(options ...Option) Option {
	return func(builder *routesBuilder) {
		builder.apply(options)
	}
}