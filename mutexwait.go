package sync

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

const (
	sleepTime = time.Millisecond
)

type MutexWait struct {
	mu     sync.Mutex
	locked atomic.Bool
}

func (m *MutexWait) TryLock() bool {
	if m.mu.TryLock() {
		m.locked.Store(true)
		return true
	}
	return false
}

func (m *MutexWait) IsLocked() bool {
	return m.locked.Load()
}

func (m *MutexWait) Lock(timeout time.Duration) bool {
	if timeout <= 0 {
		return m.TryLock()
	}

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	ticker := time.NewTicker(sleepTime)
	defer ticker.Stop()

	for {
		if m.TryLock() {
			return true
		}

		select {
		case <-timer.C:
			return false
		case <-ticker.C:
		}
	}
}

func (m *MutexWait) LockContext(ctx context.Context) bool {
	ticker := time.NewTicker(sleepTime)
	defer ticker.Stop()

	for {
		if m.TryLock() {
			return true
		}

		select {
		case <-ctx.Done():
			return false
		case <-ticker.C:
		}
	}
}

func (m *MutexWait) Unlock() {
	m.locked.Store(false)
	m.mu.Unlock()
}
