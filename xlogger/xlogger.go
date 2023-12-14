package xlogger

import (
	"io"
	"sync"
	"sync/atomic"
	"unsafe"
)

var (
	globalLogger Logger
	once         sync.Once
)

type xLogger struct {
	pool      unsafe.Pointer
	sink      *atomic.Value
	scopeName string
	logLevel  DynamicLogLevel
	isCopy    bool
	isVCalled bool
}

func newXLogger(w io.WriteCloser) Logger {
	pool := &sync.Pool{
		New: func() any {
			return &xLogger{
				isVCalled: false,
				isCopy:    true,
				sink:      &atomic.Value{},
				scopeName: "",
			}
		},
	}
	logger := &xLogger{
		sink: &atomic.Value{},
	}
	logger.setPool(pool)
	return logger
}

func (logger *xLogger) copyXLogger() *xLogger {
	_logger := logger.getPool().Get().(*xLogger)
	if _logger.getPool() == nil {
		_logger.setPool(logger.getPool())
	}
	_logger.setSink(logger.getSink())
	_logger.logLevel = logger.logLevel
	_logger.scopeName = logger.scopeName
	return _logger
}

func (logger *xLogger) releaseCopy() {
	if !logger.isCopy {
		return
	}
	logger.scopeName = ""
	logger.setSink(nil)
	logger.getPool().Put(logger)
}

func (logger *xLogger) setPool(pool *sync.Pool) {
	atomic.StorePointer(&logger.pool, unsafe.Pointer(pool))
}

func (logger *xLogger) getPool() *sync.Pool {
	return (*sync.Pool)(atomic.LoadPointer(&logger.pool))
}

func (logger *xLogger) setSink(sink Sink) {
	logger.sink.Store(sink)
}

func (logger *xLogger) getSink() Sink {
	return logger.sink.Load().(Sink)
}

func (logger *xLogger) WithName(name string) Logger {
	if logger.isCopy {
		return nil
	}
	_logger := logger.copyXLogger()
	_logger.isCopy = false
	return _logger
}

func (logger *xLogger) logLevelEnabled(lvl DynamicLogLevel) bool {
	return logger.logLevel <= lvl
}

func (logger *xLogger) V(logLevel DynamicLogLevel) Sink {
	if logger.isCopy || logger.isVCalled || !logger.logLevelEnabled(logLevel) {
		return &nilSink{}
	}
	_logger := logger.copyXLogger()
	_logger.isCopy = false
	return _logger
}

func (logger *xLogger) With(keysAndValues ...any) Sink {
	if logger.isVCalled {
		return logger.getSink().With(keysAndValues...)
	}
	return &nilSink{}
}

func (logger *xLogger) Logf(format string, args ...any) {
	defer func() {
		logger.releaseCopy()
	}()
	logger.getSink().Logf(format, args...)
}

func (logger *xLogger) Log(msg string, keysAndValues ...any) {
	defer func() {
		logger.releaseCopy()
	}()
	logger.getSink().Log(msg, keysAndValues...)
}

func (logger *xLogger) LogErr(err error, msg string, keysAndValues ...any) {
	defer func() {
		logger.releaseCopy()
	}()
	logger.getSink().LogErr(err, msg, keysAndValues...)
}
