/**------------------------------------------------------------**
 * @filename maxmind/writer.go
 * @author   jiny - caojingyin@jinycoo.com
 * @version  1.0.0
 * @date     2020/11/6 12:27
 * @desc     maxmind-writer - summary
 **------------------------------------------------------------**/

package maxmind

import (
	"crypto/md5"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jinycoo/jinygo/errors"
	"github.com/jinycoo/jinygo/log"
)

type LocalFileDatabaseWriter struct {
	filePath      string
	lockFilePath  string
	verbose       bool
	lock          *Flock
	oldHash       string
	fileWriter    io.Writer
	temporaryFile *os.File
	md5Writer     hash.Hash
}

func NewLocalFileDatabaseWriter(filePath, lockFilePath string, verbose bool) (*LocalFileDatabaseWriter, error) {
	dbWriter := &LocalFileDatabaseWriter{
		filePath:     filePath,
		lockFilePath: lockFilePath,
		verbose:      verbose,
	}

	var err error
	if dbWriter.lock, err = CreateLockFile(lockFilePath, verbose); err != nil {
		return nil, err
	}
	if err = dbWriter.createOldMD5Hash(); err != nil {
		return nil, err
	}

	temporaryFilename := fmt.Sprintf("%s.temporary", dbWriter.filePath)
	dbWriter.temporaryFile, err = os.OpenFile(
		temporaryFilename,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0644,
	)
	if err != nil {
		return nil, errors.Wrap(err, "error creating temporary file")
	}
	dbWriter.md5Writer = md5.New()
	dbWriter.fileWriter = io.MultiWriter(dbWriter.md5Writer, dbWriter.temporaryFile)

	return dbWriter, nil
}

func (writer *LocalFileDatabaseWriter) createOldMD5Hash() error {
	currentDatabaseFile, err := os.Open(writer.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			writer.oldHash = ZeroMD5
			return nil
		}
		return errors.Wrap(err, "error opening database")
	}

	defer func() {
		err := currentDatabaseFile.Close()
		if err != nil {
			log.Errorf("error(%v) closing database", err)
		}
	}()
	oldHash := md5.New()
	if _, err := io.Copy(oldHash, currentDatabaseFile); err != nil {
		return errors.Wrap(err, "error calculating database hash")
	}
	writer.oldHash = fmt.Sprintf("%x", oldHash.Sum(nil))
	if writer.verbose {
		log.Infof("Calculated MD5 sum for %s: %s", writer.filePath, writer.oldHash)
	}
	return nil
}

// Write writes to the temporary file.
func (writer *LocalFileDatabaseWriter) Write(p []byte) (int, error) {
	return writer.fileWriter.Write(p)
}

// Close closes the temporary file and releases the file lock.
func (writer *LocalFileDatabaseWriter) Close() error {
	err := writer.temporaryFile.Close()
	if err != nil {
		if perr, ok := err.(*os.PathError); !ok || perr.Err != os.ErrClosed {
			return errors.Wrap(err, "error closing temporary file")
		}
	}

	if err := os.Remove(writer.temporaryFile.Name()); err != nil && !os.IsNotExist(err) {
		return errors.Wrap(err, "error removing temporary file")
	}
	if err := writer.lock.Unlock(); err != nil {
		return errors.Wrap(err, "error releasing lock file")
	}
	return nil
}

// ValidHash checks that the temporary file's MD5 matches the given hash.
func (writer *LocalFileDatabaseWriter) ValidHash(expectedHash string) error {
	actualHash := fmt.Sprintf("%x", writer.md5Writer.Sum(nil))
	if !strings.EqualFold(actualHash, expectedHash) {
		return errors.Errorf("md5 of new database (%s) does not match expected md5 (%s)", actualHash, expectedHash)
	}
	return nil
}

// SetFileModificationTime sets the database's file access and modified times
// to the given time.
func (writer *LocalFileDatabaseWriter) SetFileModificationTime(lastModified time.Time) error {
	if err := os.Chtimes(writer.filePath, lastModified, lastModified); err != nil {
		return errors.Wrap(err, "error setting times on file")
	}
	return nil
}

// Commit renames the temporary file to the name of the database file and syncs
// the directory.
func (writer *LocalFileDatabaseWriter) Commit() error {
	if err := writer.temporaryFile.Sync(); err != nil {
		return errors.Wrap(err, "error syncing temporary file")
	}
	if err := writer.temporaryFile.Close(); err != nil {
		return errors.Wrap(err, "error closing temporary file")
	}
	if err := os.Rename(writer.temporaryFile.Name(), writer.filePath); err != nil {
		return errors.Wrap(err, "error moving database into place")
	}

	// fsync the directory. http://austingroupbugs.net/view.php?id=672
	dh, err := os.Open(filepath.Dir(writer.filePath))
	if err != nil {
		return errors.Wrap(err, "error opening database directory")
	}
	defer func() {
		if err := dh.Close(); err != nil {
			log.Fatal("Error closing directory: %+v", errors.Wrap(err, "closing directory"))
		}
	}()

	// We ignore Sync errors as they primarily happen on file systems that do
	// not support sync.
	_ = dh.Sync()
	return nil
}

// GetHash returns the hash of the current database file.
func (writer *LocalFileDatabaseWriter) GetHash() string {
	return writer.oldHash
}
