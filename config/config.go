/**------------------------------------------------------------**
 * @filename config/config.go
 * @author   jinycoo - caojingyin@jinycoo.com
 * @version  1.0.0
 * @date     2019/11/5 11:41
 * @desc     config - summary
 **------------------------------------------------------------**/
package config

import (
	"github.com/jinycoo/jinygo/errors"
	"github.com/jinycoo/jinygo/log"
	"github.com/jinycoo/jinygo/utils/file/toml"
)

var clt *Client

const DefConfigFile = "app.toml"

func TomlCfgInit(cfgPath string, v interface{}) (err error) {
	if cfgPath != "" {
		_, err = toml.DecodeFile(cfgPath, &v)
	} else {
		return remote(v)
	}
	return
}

func local(cfgPath string) (conf interface{}, err error) {
	_, err = toml.DecodeFile(cfgPath, &conf)
	return
}

func remote(v interface{}) (err error) {
	if clt, err = New(); err != nil {
		return
	}
	if err = load(v); err != nil {
		return
	}
	go func() {
		for range clt.Event() {
			if err = load(v); err != nil {
				log.Errorf("config reload error (%v)", err)
			}
		}
	}()
	return
}

func load(v interface{}) (err error) {
	var (
		s  string
		ok bool
	)
	if s, ok = clt.Toml2(); !ok {
		err = errors.New("load config center error")
	}
	if _, err = toml.Decode(s, &v); err != nil {
		err = errors.New("could not decode config")
	}
	return
}

// var (
// 	// DefaultClient default client.
// 	DefaultClient Client
// 	confPath      string
// )
//
// func init() {
// 	flag.StringVar(&confPath, "conf", "", "default config path")
// }
