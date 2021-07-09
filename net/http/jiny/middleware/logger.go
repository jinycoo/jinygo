package middleware

import (
	"net/http"
	"time"

	"github.com/jinycoo/jinygo/ctime"
	"github.com/jinycoo/jinygo/errors"
	"github.com/jinycoo/jinygo/log"
	"github.com/jinycoo/jinygo/net/http/jiny/server"
	"github.com/jinycoo/jinygo/net/metadata"
	"github.com/jinycoo/jinygo/utils/json"
)

type LogFormatterParams struct {
	StatusCode   int                    `json:"http_code"`
	ClientIP     string                 `json:"ip"`
	Method       string                 `json:"method"`
	Path         string                 `json:"path"`
	ErrorMessage string                 `json:"error_msg"`
	Keys         map[string]interface{} `json:"keys"`
	Params       string                 `json:"params"`
	TimeoutQuota float64                `json:"timeout_quota"`
	BodySize     int                    `json:"body_size"`
	Latency      int64                  `json:"latency"`
}

func Logger() server.HandlerFn {
	return func(c *server.Context) {
		start := time.Now()
		req := c.Request
		path := req.URL.Path
		raw := req.URL.RawQuery
		params := req.Form
		clientIP := metadata.String(c, metadata.RemoteIP)
		if len(clientIP) == 0 {
			clientIP = c.ClientIP()
		}
		var quota float64
		if deadline, ok := c.Deadline(); ok {
			quota = time.Until(deadline).Seconds()
		}

		c.Next()

		if req.Method != http.MethodOptions {
			logLvl := log.Info
			delay := time.Now().Sub(start)
			latency := ctime.DiffMilli(delay)
			if latency > 500 {
				logLvl = log.Warn
			}
			param := LogFormatterParams{
				Keys:         c.Keys,
				ClientIP:     clientIP,
				Method:       req.Method,
				StatusCode:   c.Writer.Status(),
				BodySize:     c.Writer.Size(),
				ErrorMessage: errors.ECause(c.Error).Message(),
				Params:       params.Encode(),
				TimeoutQuota: quota,
				Latency:      latency,
			}

			if raw != "" {
				path = path + "?" + raw
			}

			param.Path = path
			msg, _ := json.Marshal(&param)
			logLvl(string(msg))
		}
	}
}
