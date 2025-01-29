package sync_test

import (
	"context"
	"time"

	"testing"

	"github.com/alex-cos/sync"
)

func HammerMutex(muw *sync.MutexWait, loops int, cdone chan bool) {
	for i := 0; i < loops; i++ {
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Microsecond)
	defer cancel()
	for i := 0; i < loops; i++ {
		locked := muw.LockContext(ctx)
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
	for i := 0; i < 10; i++ {
		go HammerMutex(muw, 1000, c)
	}
	for i := 0; i < 10; i++ {
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
	for i := 0; i < 10; i++ {
		go HammerMutexContext(muw, 1000, c)
	}
	for i := 0; i < 10; i++ {
		b := <-c
		if !b {
			t.Fatal()
		}
	}
	close(c)
}
