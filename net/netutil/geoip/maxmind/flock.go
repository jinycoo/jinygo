/**------------------------------------------------------------**
 * @filename maxmind/flock.go
 * @author   jiny - caojingyin@baimaohui.net
 * @version  1.0.0
 * @date     2020/11/6 15:45
 * @desc     maxmind-flock - summary
 **------------------------------------------------------------**/

package maxmind

import (
	"context"
	"os"
	"runtime"
	"sync"
	"syscall"
	"time"
)

// Flock is the struct type to handle file locking. All fields are unexported,
// with access to some of the fields provided by getter methods (Path() and Locked()).
type Flock struct {
	path string
	m    sync.RWMutex
	fh   *os.File
	l    bool
	r    bool
}

// New returns a new instance of *Flock. The only parameter
// it takes is the path to the desired lockfile.
func InitFlock(path string) *Flock {
	return &Flock{path: path}
}

// NewFlock returns a new instance of *Flock. The only parameter
// it takes is the path to the desired lockfile.
//
// Deprecated: Use New instead.
func NewFlock(path string) *Flock {
	return InitFlock(path)
}

// Close is equivalent to calling Unlock.
//
// This will release the lock and close the underlying file descriptor.
// It will not remove the file from disk, that's up to your application.
func (f *Flock) Close() error {
	return f.Unlock()
}

// Path returns the path as provided in NewFlock().
func (f *Flock) Path() string {
	return f.path
}

// Locked returns the lock state (locked: true, unlocked: false).
//
// Warning: by the time you use the returned value, the state may have changed.
func (f *Flock) Locked() bool {
	f.m.RLock()
	defer f.m.RUnlock()
	return f.l
}

// RLocked returns the read lock state (locked: true, unlocked: false).
//
// Warning: by the time you use the returned value, the state may have changed.
func (f *Flock) RLocked() bool {
	f.m.RLock()
	defer f.m.RUnlock()
	return f.r
}

func (f *Flock) String() string {
	return f.path
}

// TryLockContext repeatedly tries to take an exclusive lock until one of the
// conditions is met: TryLock succeeds, TryLock fails with error, or Context
// Done channel is closed.
func (f *Flock) TryLockContext(ctx context.Context, retryDelay time.Duration) (bool, error) {
	return tryCtx(ctx, f.TryLock, retryDelay)
}

// TryRLockContext repeatedly tries to take a shared lock until one of the
// conditions is met: TryRLock succeeds, TryRLock fails with error, or Context
// Done channel is closed.
func (f *Flock) TryRLockContext(ctx context.Context, retryDelay time.Duration) (bool, error) {
	return tryCtx(ctx, f.TryRLock, retryDelay)
}

func tryCtx(ctx context.Context, fn func() (bool, error), retryDelay time.Duration) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	for {
		if ok, err := fn(); ok || err != nil {
			return ok, err
		}
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-time.After(retryDelay):
			// try again
		}
	}
}

func (f *Flock) setFh() error {
	// open a new os.File instance
	// create it if it doesn't exist, and open the file read-only.
	flags := os.O_CREATE
	if runtime.GOOS == "aix" {
		// AIX cannot preform write-lock (ie exclusive) on a
		// read-only file.
		flags |= os.O_RDWR
	} else {
		flags |= os.O_RDONLY
	}
	fh, err := os.OpenFile(f.path, flags, os.FileMode(0600))
	if err != nil {
		return err
	}

	// set the filehandle on the struct
	f.fh = fh
	return nil
}

// ensure the file handle is closed if no lock is held
func (f *Flock) ensureFhState() {
	if !f.l && !f.r && f.fh != nil {
		f.fh.Close()
		f.fh = nil
	}
}

// Lock is a blocking call to try and take an exclusive file lock. It will wait
// until it is able to obtain the exclusive file lock. It's recommended that
// TryLock() be used over this function. This function may block the ability to
// query the current Locked() or RLocked() status due to a RW-mutex lock.
//
// If we are already exclusive-locked, this function short-circuits and returns
// immediately assuming it can take the mutex lock.
//
// If the *Flock has a shared lock (RLock), this may transparently replace the
// shared lock with an exclusive lock on some UNIX-like operating systems. Be
// careful when using exclusive locks in conjunction with shared locks
// (RLock()), because calling Unlock() may accidentally release the exclusive
// lock that was once a shared lock.
func (f *Flock) Lock() error {
	return f.lock(&f.l, syscall.LOCK_EX)
}

// RLock is a blocking call to try and take a shared file lock. It will wait
// until it is able to obtain the shared file lock. It's recommended that
// TryRLock() be used over this function. This function may block the ability to
// query the current Locked() or RLocked() status due to a RW-mutex lock.
//
// If we are already shared-locked, this function short-circuits and returns
// immediately assuming it can take the mutex lock.
func (f *Flock) RLock() error {
	return f.lock(&f.r, syscall.LOCK_SH)
}

func (f *Flock) lock(locked *bool, flag int) error {
	f.m.Lock()
	defer f.m.Unlock()

	if *locked {
		return nil
	}

	if f.fh == nil {
		if err := f.setFh(); err != nil {
			return err
		}
		defer f.ensureFhState()
	}

	if err := syscall.Flock(int(f.fh.Fd()), flag); err != nil {
		shouldRetry, reopenErr := f.reopenFDOnError(err)
		if reopenErr != nil {
			return reopenErr
		}

		if !shouldRetry {
			return err
		}

		if err = syscall.Flock(int(f.fh.Fd()), flag); err != nil {
			return err
		}
	}

	*locked = true
	return nil
}

// Unlock is a function to unlock the file. This file takes a RW-mutex lock, so
// while it is running the Locked() and RLocked() functions will be blocked.
//
// This function short-circuits if we are unlocked already. If not, it calls
// syscall.LOCK_UN on the file and closes the file descriptor. It does not
// remove the file from disk. It's up to your application to do.
//
// Please note, if your shared lock became an exclusive lock this may
// unintentionally drop the exclusive lock if called by the consumer that
// believes they have a shared lock. Please see Lock() for more details.
func (f *Flock) Unlock() error {
	f.m.Lock()
	defer f.m.Unlock()

	// if we aren't locked or if the lockfile instance is nil
	// just return a nil error because we are unlocked
	if (!f.l && !f.r) || f.fh == nil {
		return nil
	}

	// mark the file as unlocked
	if err := syscall.Flock(int(f.fh.Fd()), syscall.LOCK_UN); err != nil {
		return err
	}

	f.fh.Close()

	f.l = false
	f.r = false
	f.fh = nil

	return nil
}

// TryLock is the preferred function for taking an exclusive file lock. This
// function takes an RW-mutex lock before it tries to lock the file, so there is
// the possibility that this function may block for a short time if another
// goroutine is trying to take any action.
//
// The actual file lock is non-blocking. If we are unable to get the exclusive
// file lock, the function will return false instead of waiting for the lock. If
// we get the lock, we also set the *Flock instance as being exclusive-locked.
func (f *Flock) TryLock() (bool, error) {
	return f.try(&f.l, syscall.LOCK_EX)
}

// TryRLock is the preferred function for taking a shared file lock. This
// function takes an RW-mutex lock before it tries to lock the file, so there is
// the possibility that this function may block for a short time if another
// goroutine is trying to take any action.
//
// The actual file lock is non-blocking. If we are unable to get the shared file
// lock, the function will return false instead of waiting for the lock. If we
// get the lock, we also set the *Flock instance as being share-locked.
func (f *Flock) TryRLock() (bool, error) {
	return f.try(&f.r, syscall.LOCK_SH)
}

func (f *Flock) try(locked *bool, flag int) (bool, error) {
	f.m.Lock()
	defer f.m.Unlock()

	if *locked {
		return true, nil
	}

	if f.fh == nil {
		if err := f.setFh(); err != nil {
			return false, err
		}
		defer f.ensureFhState()
	}

	var retried bool
retry:
	err := syscall.Flock(int(f.fh.Fd()), flag|syscall.LOCK_NB)

	switch err {
	case syscall.EWOULDBLOCK:
		return false, nil
	case nil:
		*locked = true
		return true, nil
	}
	if !retried {
		if shouldRetry, reopenErr := f.reopenFDOnError(err); reopenErr != nil {
			return false, reopenErr
		} else if shouldRetry {
			retried = true
			goto retry
		}
	}

	return false, err
}

// reopenFDOnError determines whether we should reopen the file handle
// in readwrite mode and try again. This comes from util-linux/sys-utils/flock.c:
//  Since Linux 3.4 (commit 55725513)
//  Probably NFSv4 where flock() is emulated by fcntl().
func (f *Flock) reopenFDOnError(err error) (bool, error) {
	if err != syscall.EIO && err != syscall.EBADF {
		return false, nil
	}
	if st, err := f.fh.Stat(); err == nil {
		// if the file is able to be read and written
		if st.Mode()&0600 == 0600 {
			f.fh.Close()
			f.fh = nil

			// reopen in read-write mode and set the filehandle
			fh, err := os.OpenFile(f.path, os.O_CREATE|os.O_RDWR, os.FileMode(0600))
			if err != nil {
				return false, err
			}
			f.fh = fh
			return true, nil
		}
	}

	return false, nil
}
