package concurrency

import (
	"context"
	"github.com/benz9527/toy-box/toys/pkg/runtime"
	"log/slog"
	"math/rand"
	"sync"
	"time"
)

type GGroup struct {
	wg sync.WaitGroup
}

func (g *GGroup) Start(fn func()) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		fn()
	}()
}

func (g *GGroup) Wait() {
	g.wg.Wait()
}

func (g *GGroup) StartWithChannel(stopC <-chan struct{}, fn func(stopC <-chan struct{})) {
	g.Start(func() {
		fn(stopC)
	})
}

func (g *GGroup) StartWithContext(ctx context.Context, fn func(ctx context.Context)) {
	g.Start(func() {
		fn(ctx)
	})
}

func ContextForChannel(parentC <-chan struct{}) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-parentC:
			cancel()
		case <-ctx.Done():
		}
	}()
	return ctx, cancel
}

var (
	neverStop <-chan struct{} = make(chan struct{})
)

type JitterFactor = float64

const (
	Factor0x JitterFactor = 0.0
	Factor1x JitterFactor = 1.0
)

type Jitter struct {
	Period        time.Duration
	Factor        JitterFactor
	Sliding       bool
	StopC         <-chan struct{}
	CrashHandlers []runtime.CrashHandler
	TraceID       string
}

func NewJitter(
	period time.Duration,
	factor JitterFactor,
	options ...func(*Jitter),
) *Jitter {
	j := &Jitter{
		Period: period,
		Factor: factor,
	}
	for _, o := range options {
		o(j)
	}
	return j
}

func WithJitterStopChannel(stopC <-chan struct{}) func(*Jitter) {
	return func(j *Jitter) {
		j.StopC = stopC
	}
}

func WithJitterCrashHandlers(handlers ...runtime.CrashHandler) func(*Jitter) {
	return func(j *Jitter) {
		j.CrashHandlers = handlers
	}
}

func WithJitterSliding(sliding bool) func(*Jitter) {
	return func(j *Jitter) {
		j.Sliding = sliding
	}
}

func WithJitterTraceID(traceID string) func(*Jitter) {
	return func(j *Jitter) {
		j.TraceID = traceID
	}
}

func (j *Jitter) extendDuration(duration time.Duration, maxFactor JitterFactor) time.Duration {
	if maxFactor <= Factor0x {
		maxFactor = Factor1x
	}
	wait := duration + time.Duration(rand.Float64()*float64(duration)*maxFactor)
	return wait
}

func (j *Jitter) Until(fn func()) {
	var (
		t          *time.Timer
		sawTimeout bool
	)

	for {
		select {
		case <-j.StopC:
			// Mitigate the case that the stopC and the timer.C were both triggered
			// (closed) at the same time.
			// Due to the fairness of the select statement,
			// the timer.C case will be picked up first.
			// So we have to double-check the stopC here.
			slog.Info("pre jitter until stopC triggered", "traceID", j.TraceID)
			return
		default:

		}

		jitterPeriod := j.Period
		if j.Factor > 0.0 {
			jitterPeriod = j.extendDuration(j.Period, j.Factor)
		}

		if !j.Sliding {
			t = ResetOrReuseTimer(t, jitterPeriod, sawTimeout)
		}

		func() {
			defer runtime.HandleCrash(false, j.CrashHandlers...)
			fn()
		}()

		if j.Sliding {
			t = ResetOrReuseTimer(t, jitterPeriod, sawTimeout)
		}

		select {
		case <-j.StopC:
			slog.Info("post jitter until stopC triggered", "traceID", j.TraceID)
			return
		case <-t.C:
			sawTimeout = true
		}
	}
}

func (j *Jitter) UntilWithContext(ctx context.Context, fn func(ctx context.Context)) {
	j.StopC = ctx.Done()
	j.Until(func() { fn(ctx) })
}

func (j *Jitter) Forever(fn func()) {
	j.Factor = Factor0x
	j.StopC = neverStop
	j.Sliding = true
	j.Until(fn)
}

func (j *Jitter) NonSlidingUntil(fn func()) {
	j.Sliding = false
	j.Factor = Factor0x
	j.Until(fn)
}

func (j *Jitter) NonSlidingUntilWithContext(ctx context.Context, fn func(context.Context)) {
	j.Sliding = false
	j.Factor = Factor0x
	j.UntilWithContext(ctx, fn)
}

func ResetOrReuseTimer(t *time.Timer, duration time.Duration, sawTimeout bool) *time.Timer {
	if t == nil {
		return time.NewTimer(duration)
	}
	if !t.Stop() && !sawTimeout {
		<-t.C
	}
	t.Reset(duration)
	return t
}
