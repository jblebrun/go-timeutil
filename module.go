package timeutil

import (
	"runtime"
	"sync"
	"time"
)

// Timer interface for swapping in testable timer
type Timer interface {
	C() <-chan time.Time
	Reset(time.Duration)
	Stop()
}

type TestTimeModule interface {
	Advance(d time.Duration)
	AwaitNewTimer()
	Now() time.Time
}

type testTimeModule struct {
	currentTime time.Time
	timers      []*testTimer
	cond        *sync.Cond
}

func (t *testTimeModule) Advance(d time.Duration) {
	t.cond.L.Lock()
	t.currentTime = t.currentTime.Add(d)
	for _, timer := range t.timers {
		if timer.fireTime != nil && timer.fireTime.Before(t.currentTime) {
			go timer.fireFunc()
			timer.fireTime = nil
			runtime.Gosched()
		}
	}
	t.cond.L.Unlock()
	runtime.Gosched()
}

func (t *testTimeModule) Now() time.Time {
	t.cond.L.Lock()
	defer t.cond.L.Unlock()
	return t.currentTime
}

func (t *testTimeModule) AwaitNewTimer() {
	t.cond.L.Lock()
	t.cond.Wait()
	t.cond.L.Unlock()
}

type timeTimer struct {
	*time.Timer
}

func (t *timeTimer) C() <-chan time.Time {
	return t.Timer.C
}
func (t *timeTimer) Reset(d time.Duration) {
	ResetTimerSafely(t.Timer, d)
}

func (t *timeTimer) Stop() {
	StopTimerSafely(t.Timer)
}

func safeTimeMake(d time.Duration, create func(d time.Duration) *time.Timer) Timer {
	var timer *time.Timer
	if d < 0 {
		timer = create(time.Hour)
		StopTimerSafely(timer)
	} else {
		timer = create(d)
	}
	return &timeTimer{timer}
}

func newTimeTimer(d time.Duration) Timer {
	return safeTimeMake(d, time.NewTimer)
}
func newTimeAfterFunc(d time.Duration, f func()) Timer {
	return safeTimeMake(d, func(d time.Duration) *time.Timer { return time.AfterFunc(d, f) })
}

type testTimer struct {
	timeModule *testTimeModule
	fireTime   *time.Time
	fireFunc   func()
	c          chan time.Time
}

func (t *testTimer) C() <-chan time.Time {
	return t.c
}

func (t *testTimer) Reset(d time.Duration) {
	t.timeModule.cond.L.Lock()
	nextFireTime := t.timeModule.currentTime.Add(d)
	t.fireTime = &nextFireTime
	t.timeModule.cond.Broadcast()
	t.timeModule.cond.L.Unlock()
}

func (t *testTimer) Stop() {
	t.timeModule.cond.L.Lock()
	t.fireTime = nil
	select {
	case <-t.c:
	default:
	}
	t.timeModule.cond.L.Unlock()
}

func newTestTimer(mod *testTimeModule, d time.Duration) Timer {
	var t *testTimer
	ch := make(chan time.Time, 1)
	t = newTestAfterFunc(mod, d, func() {
		select {
		case ch <- mod.currentTime:
		default:
		}
	}, ch).(*testTimer)
	return t
}

func newTestAfterFunc(mod *testTimeModule, d time.Duration, f func(), ch chan time.Time) Timer {
	mod.cond.L.Lock()
	defer mod.cond.L.Unlock()
	var fireTime *time.Time
	if d >= 0 {
		ft := mod.currentTime.Add(d)
		fireTime = &ft
	}

	t := &testTimer{mod, fireTime, f, ch}
	mod.timers = append(mod.timers, t)
	mod.cond.Broadcast()
	return t
}

var moduleLock sync.Mutex
var testModule *testTimeModule

func NewTest() TestTimeModule {
	moduleLock.Lock()
	defer moduleLock.Unlock()
	testModule = &testTimeModule{
		cond: sync.NewCond(&sync.Mutex{}),
	}
	return testModule
}

func ClearTest() {
	moduleLock.Lock()
	defer moduleLock.Unlock()
	testModule = nil
}

func Now() time.Time {
	moduleLock.Lock()
	defer moduleLock.Unlock()
	if testModule != nil {
		return testModule.Now()
	}
	return time.Now()
}

func NewTimer(d time.Duration) Timer {
	moduleLock.Lock()
	defer moduleLock.Unlock()
	if testModule != nil {
		return newTestTimer(testModule, d)
	}
	return newTimeTimer(d)

}

func AfterFunc(d time.Duration, f func()) Timer {
	moduleLock.Lock()
	defer moduleLock.Unlock()
	if testModule != nil {
		return newTestAfterFunc(testModule, d, f, nil)
	}
	return newTimeAfterFunc(d, f)
}
