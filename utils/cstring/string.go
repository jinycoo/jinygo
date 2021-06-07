/**------------------------------------------------------------**
 * @filename string/string.go
 * @author   jinycoo
 * @version  1.0.0
 * @date     2019-07-31 11:41
 * @desc     string - string
 **------------------------------------------------------------**/
package cstring

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"unicode"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func IsEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	if !emailRegex.MatchString(e) {
		return false
	}
	parts := strings.Split(e, "@")
	mx, err := net.LookupMX(parts[1])
	if err != nil || len(mx) == 0 {
		return false
	}
	return true
}

func LastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

func ToLower(s string) string {
	if isLower(s) {
		return s
	}

	b := make([]byte, len(s))
	for i := range b {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return BytesToString(b)
}

func isLower(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			return false
		}
	}
	return true
}

func GenEmailStar(email string) (xemail string) {
	emailBytes := []byte(email)
	el := len(emailBytes)
	if el > 10 {
		xemail = fmt.Sprintf("%s******%s", string(emailBytes[0:3]), string(emailBytes[el-4:]))
	} else {
		xemail = fmt.Sprintf("%s***%s", emailBytes[0:2], emailBytes[el-2:])
	}
	return
}

//判断是否为手机号
func IsMobileValid(mobile string) (b bool) {
	if len(mobile) < 11 {
		return false
	}
	if m, _ := regexp.MatchString("^1([38][0-9]|14[57]|5[^4])\\d{8}$", mobile); !m {
		return false
	}
	return true
}

//判断是否为手机号
func IsUrlValid(url string) (b bool) {
	if len(url) < 0 {
		return false
	}
	if b, _ := regexp.MatchString("^(?:https?:\\/\\/)?(?:[^@\\/\\n]+@)?(?:www\\.)?([^:\\/\\n]+)", url); !b {
		return false
	}
	return true
}

func ContainsPunct(body, excludes string) (rs bool) {
	var rl = []rune(body)
	for _, ch := range rl {
		if unicode.IsPunct(ch) && !strings.ContainsRune(excludes, ch) {
			rs = true
			break
		}
	}
	return
}
