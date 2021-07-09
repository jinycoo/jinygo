/**------------------------------------------------------------**
 * @filename exit/exit.go
 * @author   jinycoo
 * @version  1.0.0
 * @date     2019-07-22 10:59
 * @desc     exit - exit
 **------------------------------------------------------------**/
package exit

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jinycoo/jinygo/log"
)

func Exit(name string, close func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Infof("%s get a signal %s", name, s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Infof("%s exit", name)
			close()
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
