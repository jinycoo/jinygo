/**------------------------------------------------------------**
 * @filename mgdb/config.go
 * @author   jinycoo - caojingyin@jinycoo.com
 * @version  1.0.0
 * @date     2019/11/13 14:47
 * @desc     mgdb - mongodb config
 **------------------------------------------------------------**/
package mgdb

import (
	"jinycoo.com/jinygo/ctime"
)

type Config struct {
	Addr         string // for trace
	DSN          string // write data source name.
	Username     string
	Password     string
	Timeout      ctime.Duration
	Database     string
	IdleTimeout  ctime.Duration // connect max life time.
	QueryTimeout ctime.Duration // query sql timeout
	ExecTimeout  ctime.Duration // execute sql timeout
}
