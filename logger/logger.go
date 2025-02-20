package logger

import (
	"log"
	"os"
	"sync"

	"config"
)

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

var (
	logFile            *os.File
	loggger            *log.Logger
	logMu              sync.Mutex
	levelLabels        = []string{"debug", "info", "warning", "error", "panic"}
	logCfg             *LogConfig
	defaultCallerDepth = 2
	logPerfix          = ""
)

func SetUp(cfg *config.Config) error {
	var err error
	logCfg = &LogConfig{
		Path: cfg.LogDir,
	}
}

func Error(v ...any) {
	if logCfg.Level > ERROR {
		return
	}
	logMu.Lock()
	defer logMu.Unlock()
	setPerfix(ERROR)
	logger.PrintLn(v)
}
