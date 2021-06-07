/**------------------------------------------------------------**
 * @filename maxmind/wirter.go
 * @author   jiny - caojingyin@jinycoo.com
 * @version  1.0.0
 * @date     2020/11/6 12:25
 * @desc     maxmind-wirter - summary
 **------------------------------------------------------------**/

package maxmind

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"jinycoo.com/jinygo/errors"
)

type Writer interface {
	io.WriteCloser
	ValidHash(expectedHash string) error
	GetHash() string
	SetFileModificationTime(lastModified time.Time) error
	Commit() error
}

func CreateLockFile(lockFilePath string, verbose bool) (*Flock, error) {
	fi, err := os.Stat(filepath.Dir(lockFilePath))
	if err != nil {
		return nil, errors.Wrap(err, "database directory is not available")
	}
	if !fi.IsDir() {
		return nil, errors.New("database directory is not a directory")
	}
	lock := InitFlock(lockFilePath)
	ok, err := lock.TryLock()
	if err != nil {
		return nil, errors.Wrap(err, "error acquiring a lock")
	}
	if !ok {
		return nil, errors.Errorf("could not acquire lock on %s", lockFilePath)
	}
	if verbose {
		log.Printf("Acquired lock file lock (%s)", lockFilePath)
	}
	return lock, nil
}
