package middleware

import "github.com/ninedraft/gemax/gemax"

type Middleware func(next gemax.Handler) gemax.Handler

func Use(handler gemax.Handler, middlewares ...Middleware) gemax.Handler {
	for i := 0; i < len(middlewares); i++ {
		var wrap = middlewares[i]
		handler = wrap(handler)
	}
	return handler
}
