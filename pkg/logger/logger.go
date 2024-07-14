package logger

import (
	"path"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	lg *zap.Logger
	c  *LoggerConfig
}

type LoggerConfig struct {
	Dir []string
}

func NewLogger(c *LoggerConfig) (*Logger, error) {

	if c == nil {
		c = &LoggerConfig{
			Dir: []string{"..", "log", "undefine.log"},
		}
	}

	logLevel := zap.NewAtomicLevelAt(zapcore.DebugLevel)

	var zc = zap.Config{
		Level:             logLevel,
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "level",
			TimeKey:        "time",
			NameKey:        "name",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
		OutputPaths:      []string{path.Join(c.Dir...)},
		ErrorOutputPaths: []string{path.Join(c.Dir...)},
	}

	logger, err := zc.Build()
	if err != nil {
		return nil, err
	}
	return &Logger{lg: logger, c: c}, nil
}

func (l *Logger) Sync() {
	l.lg.Sync()
}

func (l *Logger) parseKvs(kvs ...interface{}) []zapcore.Field {
	zapFields := []zapcore.Field{}
	if len(kvs)%2 != 0 {
		zapFields = append(zapFields, zap.Any("keyError", kvs))
		return zapFields
	}
	for i := 0; i < len(kvs)-1; i += 2 {
		key := kvs[i].(string)
		zapFields = append(zapFields, zap.Any(key, kvs[i+1]))
	}
	return zapFields
}

func (l *Logger) Debug(msg string, kvs ...interface{}) {
	l.lg.Debug(msg, l.parseKvs(kvs...)...)
}

func (l *Logger) Info(msg string, kvs ...interface{}) {
	l.lg.Info(msg, l.parseKvs(kvs...)...)
}

func (l *Logger) Error(msg string, kvs ...interface{}) {
	l.lg.Error(msg, l.parseKvs(kvs...)...)
}

func (l *Logger) Panic(msg string, kvs ...interface{}) {
	l.lg.Panic(msg, l.parseKvs(kvs...)...)
}

func (l *Logger) Fatal(msg string, kvs ...interface{}) {
	l.lg.Fatal(msg, l.parseKvs(kvs...)...)
}
