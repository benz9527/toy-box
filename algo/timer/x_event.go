package timer

import (
	"sync"
	"sync/atomic"
)

type timingWheelOperation uint8

const (
	unknown timingWheelOperation = iota
	addTask
	reAddTask
	cancelTask
)

type timingWheelEvent struct {
	operation timingWheelOperation
	obj       *atomic.Value // Task or JobID
	hasSetup  bool
}

func newTimingWheelEvent(operation timingWheelOperation) *timingWheelEvent {
	event := &timingWheelEvent{
		operation: operation,
		obj:       &atomic.Value{},
	}
	return event
}

func (e *timingWheelEvent) GetOperation() timingWheelOperation {
	return e.operation
}

func (e *timingWheelEvent) GetTask() (Task, bool) {
	if e.operation != addTask && e.operation != reAddTask {
		return nil, false
	}

	obj := e.obj.Load()
	if task, ok := obj.(Task); ok {
		return task, true
	}
	return nil, false
}

func (e *timingWheelEvent) GetCancelTaskJobID() (JobID, bool) {
	if e.operation != cancelTask {
		return "", false
	}

	obj := e.obj.Load()
	if jobID, ok := obj.(JobID); ok {
		return jobID, true
	}
	return "", false
}

func (e *timingWheelEvent) CancelTaskJobID(jobID JobID) {
	if e.hasSetup {
		return
	}
	e.operation = cancelTask
	e.obj.Store(jobID)
	e.hasSetup = true
}

func (e *timingWheelEvent) AddTask(task Task) {
	if e.hasSetup {
		return
	}
	e.operation = addTask
	e.obj.Store(task)
	e.hasSetup = true
}

func (e *timingWheelEvent) ReAddTask(task Task) {
	if e.hasSetup {
		return
	}
	e.operation = reAddTask
	e.obj.Store(task)
	e.hasSetup = true
}

func (e *timingWheelEvent) clear() {
	e.operation = unknown
	e.hasSetup = false
	e.obj = &atomic.Value{}
}

type timingWheelEventsPool struct {
	pool *sync.Pool
}

func (p *timingWheelEventsPool) Get() *timingWheelEvent {
	return p.pool.Get().(*timingWheelEvent)
}

func (p *timingWheelEventsPool) Put(event *timingWheelEvent) {
	event.clear()
	p.pool.Put(event)
}

func newTimingWheelEventsPool() *timingWheelEventsPool {
	return &timingWheelEventsPool{
		pool: &sync.Pool{
			New: func() any {
				return newTimingWheelEvent(unknown)
			},
		},
	}
}
