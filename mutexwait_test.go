package sync_test

import (
	"context"
	"time"

	"testing"

	"github.com/alex-cos/sync"
)

func HammerMutex(muw *sync.MutexWait, loops int, cdone chan bool) {
	for range loops {
		locked := muw.Lock(5 * time.Microsecond)
		if !locked {
			break
		}
		defer muw.Unlock()
		locked = muw.IsLocked()
		if !locked {
			cdone <- false
			return
		}
		time.Sleep(10 * time.Microsecond)
	}
	cdone <- true
}

func HammerMutexContext(muw *sync.MutexWait, loops int, cdone chan bool) {
	for range loops {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Microsecond)
		locked := muw.LockContext(ctx)
		cancel()
		if !locked {
			break
		}
		defer muw.Unlock()
		locked = muw.IsLocked()
		if !locked {
			cdone <- false
			return
		}
		time.Sleep(10 * time.Microsecond)
	}
	cdone <- true
}

func TestMutex(t *testing.T) {
	t.Parallel()

	muw := new(sync.MutexWait)

	c := make(chan bool)
	for range 10 {
		go HammerMutex(muw, 1000, c)
	}
	for range 10 {
		b := <-c
		if !b {
			t.Fatal()
		}
	}
	close(c)
}

func TestMutexContext(t *testing.T) {
	t.Parallel()

	muw := new(sync.MutexWait)

	c := make(chan bool)
	for range 10 {
		go HammerMutexContext(muw, 1000, c)
	}
	for range 10 {
		b := <-c
		if !b {
			t.Fatal()
		}
	}
	close(c)
}
