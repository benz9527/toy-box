package concurrency

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

var (
	ErrPollConditionWaitTimeout   = errors.New("poll condition wait timeout")
	ErrPollDuplicateRun           = errors.New("poller is already running")
	ErrPollSignalGeneratorStopped = errors.New("poller signal generator is stopped")
)

type (
	PollTaskFunc func() (done bool, err error)
	PollWaitFunc func(doneC GenericWaitChannel[struct{}]) (GenericWaitChannel[struct{}], error)
)

type Poller struct {
	interval   time.Duration
	timeout    time.Duration
	isRunning  atomic.Bool
	forceStopC GenericChannel[struct{}]
}

func WithPollerForceStopChannel(stopC GenericChannel[struct{}]) func(p *Poller) {
	return func(p *Poller) {
		p.forceStopC = stopC
	}
}

func WithPollerInterval(interval time.Duration) func(p *Poller) {
	return func(p *Poller) {
		p.interval = interval
	}
}

func WithPollerTimeout(timeout time.Duration) func(p *Poller) {
	return func(p *Poller) {
		p.timeout = timeout
	}
}

func NewPoller(opts ...func(poller *Poller)) *Poller {
	p := &Poller{}
	for _, opt := range opts {
		if opt != nil {
			opt(p)
		}
	}
	if p.forceStopC == nil {
		p.forceStopC = NewSafeChannel[struct{}]()
	}
	return p
}

func (p *Poller) ForceStop() {
	_ = p.forceStopC.Close()
}

func (p *Poller) signalGen(doneC GenericWaitChannel[struct{}]) (GenericWaitChannel[struct{}], error) {
	if p.isRunning.Load() {
		return nil, ErrPollDuplicateRun
	}
	pollTaskC := NewSafeChannel[struct{}]()
	go func() {
		p.isRunning.Store(true)
		defer func() {
			_ = pollTaskC.Close()
			p.isRunning.Store(false)
			_ = p.forceStopC.Close()
		}()

		tick := time.NewTicker(p.interval)
		defer func() {
			tick.Stop()
		}()

		var afterC <-chan time.Time
		if p.timeout != 0 {
			// time.After() maybe leak if current goroutine exits early.
			timer := time.NewTimer(p.timeout)
			afterC = timer.C
			defer func() {
				timer.Stop()
			}()
		}

		for {
			select {
			case <-tick.C:
				_ = pollTaskC.Send(struct{}{}, true)
			case <-afterC:
				return
			case <-doneC.Wait():
				return
			case <-p.forceStopC.Wait():
				return
			}
		}
	}()
	return pollTaskC, nil
}

func (p *Poller) Poll(taskFn PollTaskFunc) error {
	return p.pollInternal(taskFn)
}

func (p *Poller) pollInternal(taskFn PollTaskFunc) error {
	if p.timeout == 0 || p.interval == 0 {
		return fmt.Errorf("poller timeout or interval is not set")
	}
	doneC := NewSafeChannel[struct{}]()
	defer func() {
		_ = doneC.Close()
	}()
	return waitFor(p.signalGen, taskFn, doneC)
}

// PollUntil polls until the taskFn returns true or stopC is closed.
func (p *Poller) PollUntil(taskFn PollTaskFunc, stopC GenericWaitChannel[struct{}]) error {
	p.timeout = 0
	ctx, cancel := ContextForChannel(stopC)
	defer func() {
		cancel()
	}()
	return waitFor(p.signalGen, taskFn, WaitChannelWrapper[struct{}](ctx.Done()))
}

func (p *Poller) PollInfinite(taskFn PollTaskFunc) error {
	neverDoneC := NewSafeChannel[struct{}]()
	defer func() {
		_ = neverDoneC.Close()
	}()
	return p.PollUntil(taskFn, neverDoneC)
}

func (p *Poller) PollImmediate(taskFn PollTaskFunc) error {
	return p.pollImmediateInternal(taskFn)
}

func (p *Poller) PollImmediateUntil(taskFn PollTaskFunc, stopC GenericWaitChannel[struct{}]) error {
	done, err := taskFn()
	if err != nil {
		return err
	}
	if done {
		return nil
	}
	select {
	case <-stopC.Wait():
		return ErrPollConditionWaitTimeout
	default:
		return p.PollUntil(taskFn, stopC)
	}
}

func (p *Poller) PollImmediateInfinite(taskFn PollTaskFunc) error {
	done, err := taskFn()
	if err != nil {
		return err
	}
	if done {
		return nil
	}
	return p.PollInfinite(taskFn)
}

func (p *Poller) pollImmediateInternal(taskFn PollTaskFunc) error {
	if p.interval == 0 {
		return fmt.Errorf("poller interval is not set")
	}
	done, err := taskFn()
	if err != nil {
		return err
	}
	if done {
		return nil
	}
	return p.pollInternal(taskFn)
}

// waitFor waits for the waitC to be closed or doneC to be closed.
// It also waits for each task to be done.
func waitFor(
	waitFn PollWaitFunc,
	taskFn PollTaskFunc,
	doneC GenericWaitChannel[struct{}],
) error {
	stopC := NewSafeChannel[struct{}]()
	defer func() {
		_ = stopC.Close()
	}()
	signalC, err := waitFn(stopC)
	if err != nil {
		return err
	}
	for {
		select {
		case _, open := <-signalC.Wait():
			if done, err := taskFn(); err != nil {
				return err
			} else if done {
				return nil
			}
			if !open {
				return ErrPollSignalGeneratorStopped
			}
		case <-doneC.Wait():
			return ErrPollConditionWaitTimeout
		}
	}
}
