package httpxrouter

import "net/http"

type NestingHandler interface {
	http.Handler
	Install(inner http.Handler)
}

