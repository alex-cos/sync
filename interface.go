package sync

import (
	"context"
	"time"
)

// nolint: iface
type LockerWait interface {
	TryLock() bool
	IsLocked() bool
	Lock(timeout time.Duration) bool
	LockContext(ctx context.Context) bool
	Unlock()
}
