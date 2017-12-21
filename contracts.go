package httpxrouter

import "net/http"

type nestingHandler interface {
	http.Handler
	Install(inner http.Handler)
}

