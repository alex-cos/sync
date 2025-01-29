package sync

import (
	"context"
	"sync/atomic"
	"time"
)

// A MutexWait is a mutual exclusion lock with wait timeout mechanism.
// The zero value for a Mutex is an unlocked mutex.
//
// A MutexWait must not be copied after first use.
type MutexWait struct {
	state uint32
}

// A LockerWait represents an object that can be locked and unlocked.
type LockerWait interface {
	Lock(d time.Duration) bool
	LockContext(ctx context.Context) bool
	Unlock()
	IsLocked() bool
}

const (
	mutexUnlocked = 0
	mutexLocked   = 1
	sleepTime     = time.Millisecond
)

func (thiz *MutexWait) Lock(d time.Duration) bool {
	timeout := time.After(d)
	for {
		select {
		case <-timeout:
			return false
		default:
			if thiz.tryOrWait() {
				return true
			}
		}
	}
}

func (thiz *MutexWait) LockContext(ctx context.Context) bool {
	for {
		select {
		case <-ctx.Done():
			return false
		default:
			if thiz.tryOrWait() {
				return true
			}
		}
	}
}

func (thiz *MutexWait) tryOrWait() bool {
	if atomic.CompareAndSwapUint32(&thiz.state, 0, mutexLocked) {
		return true
	}
	time.Sleep(sleepTime)
	return false
}

func (thiz *MutexWait) Unlock() {
	atomic.StoreUint32(&thiz.state, mutexUnlocked)
}

func (thiz *MutexWait) IsLocked() bool {
	return atomic.LoadUint32(&thiz.state) > 0
}
