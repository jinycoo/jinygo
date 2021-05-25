package middleware

import (
	"strconv"

	"jinycoo.com/jinygo/net/http/jiny/server"
	"jinycoo.com/jinygo/net/metadata"
	"jinycoo.com/jinygo/net/trace"
)

func Trace() server.HandlerFn {
	return func(c *server.Context) {
		t, err := trace.Extract(trace.HTTPFormat, c.Request.Header)
		if err != nil {
			var opts []trace.Option
			if ok, _ := strconv.ParseBool(trace.ETTraceDebug); ok {
				opts = append(opts, trace.EnableDebug())
			}
			t = trace.New(c.Request.URL.Path, opts...)
		}
		t.SetTitle(c.Request.URL.Path)
		t.SetTag(trace.String(trace.TagComponent, _defaultComponentName))
		t.SetTag(trace.String(trace.TagHTTPMethod, c.Request.Method))
		t.SetTag(trace.String(trace.TagHTTPURL, c.Request.URL.String()))
		t.SetTag(trace.String(trace.TagSpanKind, "server"))
		t.SetTag(trace.String("caller", metadata.String(c.Context, metadata.Caller)))
		c.Context = trace.NewContext(c.Context, t)
		c.Next()
		t.Finish(&c.Error)
	}
}
