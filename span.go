package timeutil

import "time"

type Span struct {
	start    time.Time
	duration time.Duration
}

func NewSpan(start time.Time, duration time.Duration) Span {
	return Span{start, duration}
}

func (ts Span) Start() time.Time {
	return ts.start
}

func (ts Span) Duration() time.Duration {
	return ts.duration
}

func (ts Span) End() time.Time {
	return ts.start.Add(ts.duration)
}

func (ts Span) Intersects(other Span) bool {
	if other.start.Before(ts.start) {
		return other.Intersects(ts)
	}
	return other.start.Before(ts.End())
}

func (ts Span) Union(other Span) Span {
	ret := Span{}
	if other.start.Before(ts.start) {
		ret.start = ts.start
	} else {
		ret.start = other.start
	}
	if ts.End().After(other.End()) {
		ret.duration = ts.End().Sub(ret.start)
	} else {
		ret.duration = other.End().Sub(ret.start)
	}
	return ret
}
