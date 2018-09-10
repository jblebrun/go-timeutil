package timeutil

import (
	"fmt"
	"testing"
	"time"
)

const FIRED = true

func failIf(t *testing.T, timer Timer, failIfFired bool) error {
	select {
	case <-timer.C():
		if failIfFired {
			return fmt.Errorf("Shouldn't have fired %#v", timer)
		}
	default:
		if !failIfFired {
			return fmt.Errorf("Should have fired %#v", timer)
		}
	}
	return nil
}
func TestTestTimeModule(t *testing.T) {
	m := NewTest()
	timer := NewTimer(-1)

	if e := failIf(t, timer, FIRED); e != nil {
		t.Fatal(e)
	}
	timer.Reset(1 * time.Second)

	failIf(t, timer, FIRED)

	for i := 0; i < 10; i++ {
		m.Advance(99 * time.Millisecond)
		if e := failIf(t, timer, FIRED); e != nil {
			t.Fatal(e)
		}
	}

	m.Advance(200 * time.Millisecond)

	if e := failIf(t, timer, !FIRED); e != nil {
		t.Fatal(e)
	}
}

func TestTimeSet(t *testing.T) {
}

func TestTimeTimeModule(t *testing.T) {
	t.Skip("Skipping real-time test")
	ClearTest()
	timer := NewTimer(10 * time.Millisecond)
	failIf(t, timer, FIRED)
	time.Sleep(15 * time.Millisecond)
	failIf(t, timer, !FIRED)

	timer = NewTimer(0)
	time.Sleep(15 * time.Millisecond)
	failIf(t, timer, FIRED)
}

func TestTestAfterFunc(t *testing.T) {
	m := NewTest()
	var fired bool
	AfterFunc(10*time.Second, func() {
		fired = true
	})
	for i := 0; i < 9; i++ {
		m.Advance(time.Second)
		if fired {
			t.Fatalf("NOT YET %v %v", i, Now())
		}
	}
	for i := 0; i < 5; i++ {
		m.Advance(500 * time.Millisecond)
	}
	if !fired {
		t.Fatalf("DIDN'T FIRE")
	}

}

func TestAfterFunc0(t *testing.T) {
	ClearTest()
	AfterFunc(-1, func() {
		t.Fatalf("DONT FIRE")
	})
	time.Sleep(1 * time.Millisecond)
}
