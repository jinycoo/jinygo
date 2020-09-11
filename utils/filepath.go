/**------------------------------------------------------------**
 * @filename filepath/filepath.go
 * @author   jinycoo
 * @version  1.0.0
 * @date     2019-07-31 11:40
 * @desc     filepath - file path utils
 **------------------------------------------------------------**/
package utils

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"jinycoo.com/jinygo/utils/cstring"
)

/**
 * 获取程序运行根目录位置
 */
func RootDir() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return strings.Replace(dir, "\\", "/", -1)
}

const (
	confPath = "/conf"
)

func ConfDir() string {
	//projectPath, _ := os.Getwd()
	return filepath.Join(RootDir(), confPath)
}

func joinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	appendSlash := cstring.LastChar(relativePath) == '/' && cstring.LastChar(finalPath) != '/'
	if appendSlash {
		return finalPath + "/"
	}
	return finalPath
}
