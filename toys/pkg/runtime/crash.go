package runtime

import "sync"

type CrashHandler func(any)

type crashHandlerWrapper struct {
	lock     sync.RWMutex
	handlers []CrashHandler
}

var defaultCrashHandlerWrapper *crashHandlerWrapper

func init() {
	defaultCrashHandlerWrapper = &crashHandlerWrapper{
		lock:     sync.RWMutex{},
		handlers: []CrashHandler{},
	}
}

func NoOpCrashHandler(any) {}

// SetDefaultCrashHandlers Only for testing.
func SetDefaultCrashHandlers(handlers ...CrashHandler) func() {
	defaultCrashHandlerWrapper.lock.Lock()
	originalHandlers := defaultCrashHandlerWrapper.handlers
	defaultCrashHandlerWrapper.handlers = handlers
	return func() {
		defaultCrashHandlerWrapper.handlers = originalHandlers
		defaultCrashHandlerWrapper.lock.Unlock()
	}
}

func HandleCrash(fallthroughCrash bool, handlers ...CrashHandler) {
	// Don't wrapper recover() in a function call. Otherwise,
	// the result of recover() is always nil.
	if r := recover(); r != nil {
		for _, fn := range defaultCrashHandlerWrapper.handlers {
			if fn != nil {
				fn(r)
			}
		}
		for _, fn := range handlers {
			if fn != nil {
				fn(r)
			}
		}
		if fallthroughCrash {
			panic(r)
		}
	}
}
