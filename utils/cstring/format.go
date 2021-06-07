/**------------------------------------------------------------**
 * @filename string/format.go
 * @author   jinycoo
 * @version  1.0.0
 * @date     2019-07-25 16:55
 * @desc     string - format
 **------------------------------------------------------------**/
package cstring

import "fmt"

func Sprint(template string, args ...interface{}) (message string) {
	message = template
	if message == "" && len(args) > 0 {
		message = fmt.Sprint(args...)
	} else if message != "" && len(args) > 0 {
		message = fmt.Sprintf(template, args...)
	}
	return
}

func FirstTitle(str string) string {
	if len(str) < 1 {
		return ""
	}
	strArry := []rune(str)
	if strArry[0] >= 97 && strArry[0] <= 122  {
		strArry[0] -=  32
	}
	return string(strArry)
}
