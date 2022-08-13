// Package log provides a standardised way to display logs in stdout.
package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"sync"
)

type loggerMap struct {
	Loggers map[string]*zap.Logger
}

var singleton *loggerMap
var once sync.Once

// Get returns the singleton logger instance
func Get(name string) *zap.Logger {
	once.Do(func() {
		singleton = &loggerMap{Loggers: make(map[string]*zap.Logger)}
	})

	logger, ok := singleton.Loggers[name]
	if !ok {
		var err error
		config := zap.Config{
			Encoding:    "console",
			Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
			OutputPaths: []string{"stdout"},
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:     "time",
				EncodeTime:  zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
				LevelKey:    "level",
				EncodeLevel: zapcore.CapitalColorLevelEncoder,
				MessageKey:  "message",
				NameKey:     "name",
				EncodeName:  zapcore.FullNameEncoder,
			},
		}

		if os.Getenv("ENV") == "prod" {
			// Only log warnings
			config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
		}

		logger, err = config.Build()

		if err != nil {
			log.Fatalln("Error loading logger:", err)
		}

		singleton.Loggers[name] = logger.Named(name)
		logger = singleton.Loggers[name]
	}
	return logger
}
