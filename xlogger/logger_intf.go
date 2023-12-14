package xlogger

type DynamicLogLevel uint8

const (
	DebugLevel DynamicLogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

type XDumpOutputType uint8

const (
	DumpAsPlainText XDumpOutputType = iota
	DumpAsJson
	DumpAsYaml
	DumpAsTable
	DumpAsXml
)

type XDumpFunc func(obj any) (XDumpOutputType, string)

type XCommonLogDumper interface {
	// Dump is a helper function to dump object to log,
	// only works in trace/debug mode
	Dump(obj any, fn XDumpFunc)
}

type Sink interface {
	V(logLevel DynamicLogLevel) Sink
	With(keysAndValues ...any) Sink
	Logf(format string, args ...any)
	Log(msg string, keysAndValues ...any)
	LogErr(err error, msg string, keysAndValues ...any)
}

type Logger interface {
	//Sink
	//XCommonLogger
	//WithName(scopeName string) Logger
	//WithKeysAndValues(keysAndValues ...any) Logger
}
