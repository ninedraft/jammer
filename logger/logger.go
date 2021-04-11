package logger

import (
	"context"
	"log"
	"time"

	"github.com/ninedraft/gemax/gemax"
	"github.com/ninedraft/gemax/gemax/status"
)

func Log(next gemax.Handler) gemax.Handler {
	return func(ctx context.Context, rw gemax.ResponseWriter, req gemax.IncomingRequest) {
		var start = time.Now()
		var inter = &interceptor{rw: rw}
		next(ctx, inter, req)
		log.Printf("request %s -> %s in %s", req.URL(), inter.status, time.Since(start))
	}
}

type interceptor struct {
	status status.Code
	rw     gemax.ResponseWriter
}

func (i *interceptor) WriteStatus(code status.Code, meta string) {
	i.status = code
	i.rw.WriteStatus(code, meta)
}

func (i *interceptor) Write(p []byte) (n int, err error) {
	if i.status == 0 {
		i.status = status.Success
	}
	return i.rw.Write(p)
}

func (i *interceptor) Close() error {
	if i.status == 0 {
		i.status = status.Success
	}
	return i.rw.Close()
}
