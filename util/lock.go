package util

import (
	"sync"
	"time"
)

// MutexWrapper is an interface implemented by Mutex and RWMutex.
type MutexWrapper interface {
	// Name returns the name of this mutex wrapper.
	Name() string
	// Timeout returns the timeout value of this mutex wrapper.
	Timeout() time.Duration
}

// Mutex is the wrapper around sync.Mutex.
type Mutex struct {
	m *sync.Mutex
	// Name represents the name of this lock.
	name string
	// Timeout configures the timeout duration to log error messages when
	// acquiring the lock. Defaults to 0, which means no message will be
	// logged.
	timeout time.Duration
}

// NewMutex constructs a new mutex.
func NewMutex(name string, timeout time.Duration) *Mutex {
	return &Mutex{
		m:       &sync.Mutex{},
		name:    name,
		timeout: timeout,
	}
}

// Name implements MutexWrapper.Name().
func (m *Mutex) Name() string {
	return m.name
}

// Timeout implements Mutex.Wrapper.Timeout().
func (m *Mutex) Timeout() time.Duration {
	return m.timeout
}

const (
	OpMutexLock    = "Mutex.Lock"
	OpRWMutexLock  = "RWMutex.Lock"
	OpRWMutexRLock = "RWMutex.RLock"
)

// Lock wraps sync.Mutex.Lock().
func (m *Mutex) Lock() {
	// Mutex with Timeout==0 equals sync.Mutex.
	if m.timeout == time.Duration(0) {
		m.m.Lock()
		return
	}
	// Make a channel with one-element buffer to receive acquire ok message.
	got := make(chan bool, 1)
	go func() {
		select {
		case <-time.After(m.timeout):
			logError(m, OpMutexLock)
		case <-got:
		}
	}()
	m.m.Lock()
	got <- true
}

// Unlock wraps sync.Mutex.Unlock().
func (m *Mutex) Unlock() {
	m.m.Unlock()
}

// RWMutex is the wrapper around sync.RWMutex
type RWMutex struct {
	m *sync.RWMutex
	// Name represents the name of this lock.
	name string
	// Timeout configures the timeout duration to log error messages when
	// acquiring the lock. Defaults to 0, which means no message will be
	// logged.
	timeout time.Duration
}

// NewRWMutex constructs a new rwmutex.
func NewRWMutex(name string, timeout time.Duration) *RWMutex {
	return &RWMutex{
		m:       &sync.RWMutex{},
		name:    name,
		timeout: timeout,
	}
}

// Name implements MutexWrapper.Name().
func (m *RWMutex) Name() string {
	return m.name
}

// Timeout implements Mutex.Wrapper.Timeout().
func (m *RWMutex) Timeout() time.Duration {
	return m.timeout
}

// Lock wraps sync.RWMutex.Lock().
func (m *RWMutex) Lock() {
	// Mutex with Timeout==0 equals sync.Mutex.
	if m.timeout == time.Duration(0) {
		m.m.Lock()
		return
	}
	// Make a channel with one-element buffer to receive acquire ok message.
	got := make(chan bool, 1)
	go func() {
		select {
		case <-time.After(m.timeout):
			logError(m, OpRWMutexLock)
		case <-got:
		}
	}()
	m.m.Lock()
	got <- true
}

// Unlock wraps sync.RWMutex.Unlock().
func (m *RWMutex) Unlock() {
	m.m.Unlock()
}

// RLock wraps sync.RWMutex.RLock().
func (m *RWMutex) RLock() {
	// Mutex with timeout==0 equals sync.Mutex.
	if m.timeout == time.Duration(0) {
		m.m.RLock()
		return
	}
	// Make a channel with one-element buffer to receive acquire ok message.
	got := make(chan bool, 1)
	go func() {
		select {
		case <-time.After(m.timeout):
			logError(m, OpRWMutexRLock)
		case <-got:
		}
	}()
	m.m.RLock()
	got <- true
}

// RUnlock wraps sync.RWMutex.RUnlock().
func (m *RWMutex) RUnlock() {
	m.m.RUnlock()
}
