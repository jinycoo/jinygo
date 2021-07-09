package jiny

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/jinycoo/jinygo/errors"
	"github.com/jinycoo/jinygo/log"
	"github.com/jinycoo/jinygo/net/http/jiny/binding"
	"github.com/jinycoo/jinygo/net/http/jiny/server"
	"github.com/jinycoo/jinygo/net/netutil/signature"
)

func SignValid() server.HandlerFn {
	return func(c *server.Context) {
		if c.Request.Method == http.MethodGet {
			var signFlag bool
			u, _ := url.Parse(c.Request.RequestURI)
			for _, p := range conf.SignPaths {
				if p == u.Path {
					signFlag = true
					break
				}
			}
			if signFlag {
				appID := c.DefaultQuery("app_id", "")
				sign := c.DefaultQuery("sign", "")
				sign = strings.Replace(sign, " ", "+", -1)

				if len(appID) == 0 || len(sign) == 0 {
					log.Warnf("client signature is invalid, appID=%v, sign=%v", appID, sign)
					c.AbortWithStatus(http.StatusUnauthorized)
				}

				if appID == conf.Sign.AppID {
					pubKeys := conf.Sign.PubKeys

					var (
						content string
						err     error
					)
					switch {
					case c.Request.Method == "GET":
						content, err = convertURLValToSignString(c.Request.RequestURI)
					case c.ContentType() == "application/json":
						//bodyDataStrTmp := getBodyData(c)
						//if !json.Valid([]byte(bodyDataStrTmp)) { //校验是否为json格式
						//	log.Infof("The request param error.bodyData:", bodyDataStrTmp)
						//	c.AbortWithStatus(http.StatusForbidden)
						//}
						//_ = json.Unmarshal([]byte(bodyDataStrTmp), &bodyDataKV)
						content, err = convertBodyToSignJSON(c)
					default:
						content, err = convertURLValToSignString("")
					}

					if err != nil {
						log.Errorf("VerifySign failed, appID=%v, content=%v, sign=%v, err=%v", appID, content, sign, err)
						c.AbortWithStatus(http.StatusForbidden)
						return
					}
					content = strings.Replace(content, " ", "+", -1)
					err = signature.VerifySign(content, sign, pubKeys)
					if err != nil {
						log.Warnf("VerifySign failed, appID=%v, content=%v, sign=%v, err=%v", appID, content, sign, err)
						c.AbortWithStatus(http.StatusForbidden)
						return
					}
					return
				} else {
					log.Warn("appID invalid")
					c.AbortWithStatus(http.StatusForbidden)
				}
			}
		}
		c.Next()
	}
}

func convertURLValToSignString(uri string) (content string, err error) {
	u, err := url.Parse(uri)
	values, _ := url.ParseQuery(u.RawQuery)
	if len(values) == 0 {
		err = errors.ParamsErr
		return
	}
	sc := make(map[string]interface{})
	var formDataKeys = make([]string, 0)
	for k, v := range values {
		var val = strings.Join(v, "")
		if len(val) == 0 || k == "app_id" || k == "sign" {
			continue
		}
		formDataKeys = append(formDataKeys, k)
		sc[k] = val
	}

	sort.Strings(formDataKeys)
	for _, k := range formDataKeys {
		content += k
		content += fmt.Sprintf("%v", sc[k])
	}
	return
}

func convertBodyToSignJSON(c *server.Context) (content string, err error) {
	data, err := c.GetRawData()
	if err != nil {
		return
	}
	sc := make(map[string]interface{})
	err = json.Unmarshal(data, &sc)
	if err != nil {
		return
	}
	var formDataKeys = make([]string, 0)
	for k, v := range sc {
		if len(strings.TrimSpace(fmt.Sprintf("%v", v))) == 0 {
			continue
		}
		formDataKeys = append(formDataKeys, k)
	}
	sort.Strings(formDataKeys)
	for _, k := range formDataKeys {
		content += k
		content += fmt.Sprintf("%v", sc[k])
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(data))
	return
}

/*
	功能：获取所有请求过来的参数
	author: haima
	time:2021.01.28
*/
func SignValid1(c *server.Context) {
	var method, requestType = c.Request.Method, c.ContentType()
	quireDateStr := ""
	var bodyDataKV = map[string]interface{}{}
	switch {
	case method == "GET": //测试通过
		//?name=lisi&age=30&gender=男
		qparam := c.Request.RequestURI
		//解析这个 URL 并确保解析没有出错。
		u, err := url.Parse(qparam)
		if err != nil {
			log.Infof("The request param error. err:", err)
			c.AbortWithStatus(http.StatusForbidden)
		}
		//直接访问 scheme。
		formData, _ := url.ParseQuery(u.RawQuery)
		for k, v := range formData {
			bodyDataKV[k] = v[0]
		}
	case requestType == binding.MIMEJSON: //测试通过
		//{
		//	"name":"lisi",
		//	"age":18,
		//	"denger":"男"
		//}
		bodyDataStrTmp := getBodyData(c)
		if !json.Valid([]byte(bodyDataStrTmp)) { //校验是否为json格式
			log.Infof("The request param error.bodyData:", bodyDataStrTmp)
			c.AbortWithStatus(http.StatusForbidden)
		}
		_ = json.Unmarshal([]byte(bodyDataStrTmp), &bodyDataKV)
	case requestType == binding.MIMEPOSTForm: //测试通过
		qdataTmp := getBodyData(c) //qdataTmp = name=haima&age=12
		qdataSlice := strings.Split(qdataTmp, "&")
		for _, v := range qdataSlice {
			vKVSlice := strings.Split(v, "=")
			bodyDataKV[vKVSlice[0]] = vKVSlice[1]
		}
	case requestType == binding.MIMEMultipartPOSTForm: //测试通过
		formData, _ := c.MultipartForm()
		for k, v := range formData.Value {
			bodyDataKV[k] = v[0]
		}
	default:
		log.Errorf("The request ContentType error. method(%s) requestType(%s)", method, requestType)
		c.AbortWithStatus(http.StatusForbidden)
	}
	quireDateStr = getPOSTFormSignStr(bodyDataKV)

	fmt.Printf("method(%s) requestType(%s) dataStr(%s) \n", method, requestType, quireDateStr)
	c.Next()
}

//排序拼接前端传过来的参数
func getPOSTFormSignStr(bodyData map[string]interface{}) string {
	var formDataDataKV = map[string]interface{}{}
	var formDataKeys []string
	for k, v := range bodyData {
		if len(strings.Replace(fmt.Sprintf("%v", v), " ", "", -1)) <= 0 {
			continue
		}
		formDataKeys = append(formDataKeys, k)
		formDataDataKV[k] = v
	}
	sort.Strings(formDataKeys)
	var bodyDataStr = ""
	for _, k := range formDataKeys {
		bodyDataStr += k
		bodyDataStr += fmt.Sprintf("%v", formDataDataKV[k])
	}
	return bodyDataStr
}

// 获取POST请求 body里的数据
func getBodyData(c *server.Context) string {
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = io.ReadAll(c.Request.Body)
	}
	// 把刚刚读出来的再写进去
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	body, _ := url.QueryUnescape(string(bodyBytes))
	return body
}
