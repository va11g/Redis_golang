package logger

type LogLevel int
type LogConfig struct {
	Path  string
	Name  string
	Level LogLevel
}

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	PANIC
)

func Error(v ...any) {
	if logcfg.Level > ERROR {
		return
	}
	logMu.Lock()
	defer logMu.Unlock()
	setPerfix(ERROR)
	logger.PrintLn(v)
}
