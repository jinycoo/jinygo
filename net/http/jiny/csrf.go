/**------------------------------------------------------------**
 * @filename jiny/csrf.go
 * @author   jinycoo
 * @version  1.0.0
 * @date     2019-07-24 10:11
 * @desc     jiny - csrf header validate
 **------------------------------------------------------------**/
package jiny

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"jinycoo.com/jinygo/net/http/jiny/server"
)

var validations = []func(*url.URL) bool{}

func matchHostSuffix(suffix string) func(*url.URL) bool {
	return func(uri *url.URL) bool {
		return strings.HasSuffix(strings.ToLower(uri.Host), suffix)
	}
}

func matchPattern(pattern *regexp.Regexp) func(*url.URL) bool {
	return func(uri *url.URL) bool {
		return pattern.MatchString(strings.ToLower(uri.String()))
	}
}

// addHostSuffix add host suffix into validations
func addHostSuffix(suffix string) {
	validations = append(validations, matchHostSuffix(suffix))
}

// addPattern add referer pattern into validations
func addPattern(pattern string) {
	validations = append(validations, matchPattern(regexp.MustCompile(pattern)))
}

func CSRF() server.HandlerFn {
	return func(c *server.Context) {
		referer := c.GetHeader("Referer")
		u, _ := url.Parse(c.Request.RequestURI)
		for _, p := range conf.AllowPatterns {
			if p == u.Path && len(referer) == 0 {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		}
		//params := c.Request.Form
		//cross := (params.Get("callback") != "" && params.Get("jsonp") == "jsonp") || (params.Get("cross_domain") != "")
		//if len(referer) == 0 {
		//	//if !cross {
		//	//	return
		//	//}
		//	log.Info("The request's Referer header is empty.")
		//	c.AbortWithStatus(http.StatusForbidden)
		//	return
		//}
		//illegal := true
		//if uri, err := url.Parse(referer); err == nil && uri.Host != "" {
		//	for _, validate := range validations {
		//		if validate(uri) {
		//			illegal = false
		//			break
		//		}
		//	}
		//}
		//if illegal {
		//	log.Infof("The request's Referer header `%s` does not match any of allowed referers.", referer)
		//	c.AbortWithStatus(http.StatusForbidden)
		//	return
		//}
	}
}
