package sync_test

import (
	"context"
	stdsync "sync"
	"sync/atomic"
	"time"

	"testing"

	"github.com/alex-cos/sync"
)

func HammerMutex(muw *sync.MutexWait, loops int, failed *atomic.Bool) {
	for range loops {
		locked := muw.Lock(50 * time.Millisecond)
		if !locked {
			break
		}
		locked = muw.IsLocked()
		if !locked {
			failed.Store(true)
			return
		}
		time.Sleep(10 * time.Microsecond)
		muw.Unlock()
	}
}

func HammerMutexContext(muw *sync.MutexWait, loops int, failed *atomic.Bool) {
	for range loops {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		locked := muw.LockContext(ctx)
		cancel()
		if !locked {
			break
		}
		locked = muw.IsLocked()
		if !locked {
			failed.Store(true)
			return
		}
		time.Sleep(10 * time.Microsecond)
		muw.Unlock()
	}
}

func TestMutex(t *testing.T) {
	t.Parallel()

	muw := new(sync.MutexWait)

	var wg stdsync.WaitGroup
	var failed atomic.Bool

	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			HammerMutex(muw, 1000, &failed)
		}()
	}
	wg.Wait()

	if failed.Load() {
		t.Fatal("mutex hammer test failed")
	}
}

func TestMutexContext(t *testing.T) {
	t.Parallel()

	muw := new(sync.MutexWait)

	var wg stdsync.WaitGroup
	var failed atomic.Bool

	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			HammerMutexContext(muw, 1000, &failed)
		}()
	}
	wg.Wait()

	if failed.Load() {
		t.Fatal("mutex context hammer test failed")
	}
}

func TestDoubleUnlock(t *testing.T) {
	t.Parallel()

	muw := new(sync.MutexWait)

	muw.LockInfinite()
	muw.Unlock()
	muw.Unlock()

	if muw.IsLocked() {
		t.Fatal("mutex should not be locked after double unlock")
	}
}

func TestContextCancelled(t *testing.T) {
	t.Parallel()

	muw := new(sync.MutexWait)
	muw.LockInfinite()

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan bool)
	go func() {
		result := muw.LockContext(ctx)
		done <- result
	}()

	time.Sleep(10 * time.Millisecond)
	cancel()

	select {
	case result := <-done:
		if result {
			t.Fatal("LockContext should return false when context is cancelled")
		}
	case <-time.After(time.Second):
		t.Fatal("LockContext did not return after context cancellation")
	}

	muw.Unlock()
}

func TestTimeoutZero(t *testing.T) {
	t.Parallel()

	muw := new(sync.MutexWait)

	result := muw.Lock(0)
	if !result {
		t.Fatal("Lock(0) should succeed on unlocked mutex")
	}
	if !muw.IsLocked() {
		t.Fatal("mutex should be locked")
	}

	result = muw.Lock(0)
	if result {
		t.Fatal("Lock(0) should fail on already locked mutex")
	}

	muw.Unlock()
}

func TestNegativeTimeout(t *testing.T) {
	t.Parallel()

	muw := new(sync.MutexWait)

	result := muw.Lock(-1)
	if !result {
		t.Fatal("Lock(-1) should behave like TryLock and succeed on unlocked mutex")
	}

	muw.Unlock()
}
