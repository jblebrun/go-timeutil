package timeutil

import (
	"time"
)

type Duration time.Duration

func (d *Duration) UnmarshalText(text []byte) error {
	dd, err := time.ParseDuration(string(text))
	*d = Duration(dd)
	return err
}

// Shamelessly lifted from
// https://play.golang.org/p/QHocTHl8iR
// ref:
// http://grokbase.com/t/gg/golang-nuts/1492epp0qb/go-nuts-how-to-round-a-duration
func RoundDuration(d, r time.Duration) time.Duration {
	neg := d < 0
	if neg {
		d = -d
	}
	if m := d % r; m+m < r {
		d = d - m
	} else {
		d = d + r - m
	}
	if neg {
		return -d
	}
	return d
}
