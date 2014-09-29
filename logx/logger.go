package logx

type Logger interface {
	// Init(opts Options) Logger
	Log(level string, message string, metadata interface{})
}
