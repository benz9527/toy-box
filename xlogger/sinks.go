package xlogger

var (
	_ Sink = (*nilSink)(nil)
)

type nilSink struct{}

func (n *nilSink) V(logLevel DynamicLogLevel) Sink                    { return n }
func (n *nilSink) With(keysAndValues ...any) Sink                     { return n }
func (n *nilSink) Logf(format string, args ...any)                    {}
func (n *nilSink) Log(msg string, keysAndValues ...any)               {}
func (n *nilSink) LogErr(err error, msg string, keysAndValues ...any) {}
func (n *nilSink) LogLevelEnabled(lvl DynamicLogLevel) bool           { return false }
