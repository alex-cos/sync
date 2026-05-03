# sync

[![Go Version](https://img.shields.io/badge/Go-1.23%2B-blue)](https://go.dev/)
[![Test Status](https://github.com/alex-cos/sync/actions/workflows/test.yml/badge.svg)](https://github.com/alex-cos/sync/actions/workflows/test.yml)
[![Lint Status](https://github.com/alex-cos/sync/actions/workflows/lint.yml/badge.svg)](https://github.com/alex-cos/sync/actions/workflows/lint.yml)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/alex-cos/sync)](https://goreportcard.com/report/github.com/alex-cos/sync)

A Go package providing `MutexWait`, a mutex wrapper with timeout, context support, and non-blocking lock attempts.

## Features

- **TryLock** — non-blocking lock attempt
- **Lock** — lock with a configurable timeout
- **LockContext** — lock with context cancellation support
- **LockInfinite** — blocking lock without timeout
- **IsLocked** — check if the mutex is currently held
- **Double-Unlock safe** — calling `Unlock` on an already unlocked mutex is a no-op

## Installation

```bash
go get github.com/alex-cos/sync
```

## Usage

### Basic lock with timeout

```go
muw := &sync.MutexWait{}

if muw.Lock(500 * time.Millisecond) {
    defer muw.Unlock()
    // critical section
} else {
    // timed out
}
```

### Context-aware lock

```go
muw := &sync.MutexWait{}

ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()

if muw.LockContext(ctx) {
    defer muw.Unlock()
    // critical section
} else {
    // context cancelled or timed out
}
```

### Non-blocking attempt

```go
muw := &sync.MutexWait{}

if muw.TryLock() {
    defer muw.Unlock()
    // acquired immediately
} else {
    // already locked, do something else
}
```

### Check lock state

```go
muw := &sync.MutexWait{}

muw.LockInfinite()
fmt.Println(muw.IsLocked()) // true

muw.Unlock()
fmt.Println(muw.IsLocked()) // false
```

### Infinite lock

```go
muw := &sync.MutexWait{}

muw.LockInfinite()
defer muw.Unlock()
// critical section — blocks until acquired
```

## Interface

The `LockerWait` interface is also exported for dependency injection and mocking:

```go
type LockerWait interface {
    TryLock() bool
    IsLocked() bool
    LockInfinite()
    Lock(timeout time.Duration) bool
    LockContext(ctx context.Context) bool
    Unlock()
}
```
