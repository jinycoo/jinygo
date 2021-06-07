/**------------------------------------------------------------**
 * @filename jiny/jiny.go
 * @author   jinycoo
 * @version  1.0.0
 * @date     2019-07-24 09:41
 * @desc     jiny - http server
 **------------------------------------------------------------**/
package jiny

import (
	"net/http"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"jinycoo.com/jinygo/errors"
	"jinycoo.com/jinygo/net/http/jiny/middleware"
	"jinycoo.com/jinygo/net/http/jiny/server"
)

var (
	once           sync.Once
	internalEngine *server.Engine
	conf           *Config
)
/**
* init
* @method init
* @return
 */
func init() {
	all := make([]string, 0)
	conf = &Config{
		Port:          ":80",
		AllowHosts:    all,
		AllowPatterns: all,
		Headers: map[string]string{
			"access-control-allow-origin":      "*",
			"access-control-allow-headers":     "*",
			"access-control-allow-methods":     "GET, OPTIONS, POST, PUT, DELETE",
			"access-control-allow-credentials": "true",
		},
	}
}

func Init(cfg *Config) {
	if cfg != nil {
		if strings.HasPrefix(cfg.Port, ":") {
			conf.Port = cfg.Port
		}
		conf.SignActive = cfg.SignActive
		conf.AllowHosts = cfg.AllowHosts
		conf.AllowPatterns = cfg.AllowPatterns
		conf.Sign = cfg.Sign
		conf.SignPaths = cfg.SignPaths
		conf.Expiry = cfg.Expiry
		conf.SigningKey = cfg.SigningKey
		for hk, hv := range cfg.Headers {
			headKey := strings.ToLower(hk)
			if _, ok := conf.Headers[headKey]; ok {
				conf.Headers[headKey] = hv
			}
		}
		for _, r := range conf.AllowHosts {
			addHostSuffix(r)
		}
		for _, p := range conf.AllowPatterns {
			addPattern(p)
		}
	}
	once.Do(func() {
		internalEngine = server.New()
		internalEngine.Use(middleware.NewMiddlewares()...)
		internalEngine.Use(CSRF())
		if conf.SignActive {
			internalEngine.Use(SignValid()) //获取所有请求过来的参数
		}
		internalEngine.Use(func(c *server.Context) {
			method := c.Request.Method
			for k, v := range conf.Headers {
				c.Header(k, v)
			}

			//origin := c.Request.Header.Get("Origin")
			//isAccess, _ := regexp.MatchString(`^(http|https)://.*.127.0.0.1(:\d+)?`, origin)
			//if isAccess {
			c.Set("Content-Type", "application/json")
			//}
			//放行所有OPTIONS方法
			if method == "OPTIONS" {
				c.JSON("OPTIONS", nil)
				c.AbortWithStatus(200)
			}
			c.Next()
		})
		internalEngine.NoRoute(func(c *server.Context) {
			err := errors.NothingFound
			c.JSON(http.StatusNotFound, err)
		})
		internalEngine.Use(middleware.Prom(nil))
		//internalEngine.addRoute("GET", "/metrics", HandlersChain{monitor()})
		//internalEngine.addRoute("GET", "/metadata", HandlersChain{engine.metadata()})

		//debug.StartPprof()
	})
	//return internalEngine
}

// Group creates a new router group. You should add all the routes that have common middlewares or the same path prefix.
// For example, all the routes that use a common middleware for authorization could be grouped.
func Group(relativePath string, handlers ...server.HandlerFn) *server.RouterGroup {
	return internalEngine.Group(relativePath, handlers...)
}

func AuthJwtGroup(relativePath string) *server.RouterGroup {
	return internalEngine.Group(relativePath, JwtAuth())
}

// Routes returns a slice of registered routes.
func Routes() server.RoutesInfo {
	return internalEngine.Routes()
}

// Ping is used to set the general HTTP ping handler.
func Ping(handler server.HandlerFn) {
	internalEngine.GET("/ping", handler)
}

func Index(handler server.HandlerFn) {
	internalEngine.GET("/", handler)
}

func PromMonitor() {
	internalEngine.GET("/metrics", middleware.PromHandler(promhttp.Handler()))
}

func Run() (err error) {
	return internalEngine.Run(conf.Port)
}
