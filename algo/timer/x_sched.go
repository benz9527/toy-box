package timer

import "time"

type xScheduler struct {
	intervals    []time.Duration
	currentIndex int
	isFinite     bool
}

var (
	_ Scheduler = (*xScheduler)(nil)
)

func NewFiniteScheduler(intervals ...time.Duration) Scheduler {
	if len(intervals) == 0 {
		return nil
	}
	for _, interval := range intervals {
		if interval.Milliseconds() <= 0 {
			return nil
		}
	}
	return &xScheduler{
		isFinite:     true,
		intervals:    intervals,
		currentIndex: 0,
	}
}

func NewInfiniteScheduler(intervals ...time.Duration) Scheduler {
	if len(intervals) == 0 {
		return nil
	}
	for _, interval := range intervals {
		if interval.Milliseconds() <= 0 {
			return nil
		}
	}
	return &xScheduler{
		intervals:    intervals,
		currentIndex: 0,
	}
}

func (x *xScheduler) next(beginMs int64) (nextExpiredMs int64) {
	beginTime := MillisToUTCTime(beginMs)
	if beginTime.IsZero() || len(x.intervals) == 0 {
		return -1
	}

	if x.currentIndex >= len(x.intervals) {
		if x.isFinite {
			return -1
		}
		x.currentIndex = 0
	}
	if x.intervals[x.currentIndex].Milliseconds() <= 0 {
		return -1
	}
	next := beginTime.Add(x.intervals[x.currentIndex])
	x.currentIndex++
	return next.UnixMilli()
}

func (x *xScheduler) GetRestLoopCount() int64 {
	if x.isFinite {
		return int64(len(x.intervals) - x.currentIndex)
	}
	return -1
}

func MillisToUTCTime(millis int64) time.Time {
	return time.Unix(0, millis*int64(time.Millisecond)).UTC()
}
