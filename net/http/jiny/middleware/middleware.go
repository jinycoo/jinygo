package middleware

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"

	"github.com/jinycoo/jinygo/log"
	"github.com/jinycoo/jinygo/net/http/jiny/server"
)

const _defaultComponentName = "net/http"

func NewMiddlewares() []server.HandlerFn {
	return []server.HandlerFn{Logger(), Recovery(), Trace()}
}

func Recovery() server.HandlerFn {
	return func(c *server.Context) {
		defer func() {
			var rawReq []byte
			if err := recover(); err != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				if c.Request != nil {
					rawReq, _ = httputil.DumpRequest(c.Request, false)
				}
				pl := fmt.Sprintf("[Recovery] http call panic: %s\n%v\n%s\n", string(rawReq), err, buf)
				_, err = fmt.Fprintf(os.Stderr, pl)
				log.ZError(pl)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
