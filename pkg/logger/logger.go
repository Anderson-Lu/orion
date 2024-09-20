package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	lg *zap.Logger
	c  *LoggerConfig
}

type LoggerConfig struct {
	Path              []string
	LogFileMaxSizeMB  int
	LogFileMaxBackups int
	LogMaxAgeDays     int
	LogCompress       bool
	LogLevel          string
}

func (l *LoggerConfig) check() {
	if l.LogFileMaxSizeMB <= 0 {
		l.LogFileMaxSizeMB = 128
	}
	if l.LogFileMaxBackups <= 0 {
		l.LogFileMaxBackups = 10
	}
	if l.LogMaxAgeDays <= 0 {
		l.LogMaxAgeDays = 30
	}
	if l.LogLevel == "" {
		l.LogLevel = "debug"
	}
}

func NewLogger(c *LoggerConfig) (*Logger, error) {

	if c == nil {
		c = &LoggerConfig{
			Path: []string{"..", "log", "undefine.log"},
		}
	}
	c.check()

	lc := &Logger{c: c}
	if err := lc.checkOutputPath(); err != nil {
		return nil, err
	}

	lvs := map[string]zapcore.Level{
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
	}
	cLv, ok := lvs[c.LogLevel]
	if !ok {
		cLv = lvs["info"]
	}

	ec := zapcore.EncoderConfig{
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
	}
	encoder := zapcore.NewJSONEncoder(ec)

	lj := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(c.Path...),
		MaxSize:    c.LogFileMaxSizeMB,
		MaxBackups: c.LogFileMaxBackups,
		MaxAge:     c.LogMaxAgeDays,
		Compress:   c.LogCompress,
	})

	lc.lg = zap.New(zapcore.NewCore(encoder, lj, cLv), zap.AddCaller(), zap.AddCallerSkip(1))
	return lc, nil
}

func (l *Logger) checkOutputPath() error {
	if len(l.c.Path) < 1 {
		return fmt.Errorf("invalid logger output path: %+v", l.c.Path)
	}
	p := filepath.Join(l.c.Path[0 : len(l.c.Path)-1]...)
	_, err := os.Stat(p)
	switch err {
	case nil:
		return nil
	default:
		if os.IsExist(err) {
			return nil
		}
		return os.Mkdir(p, os.ModePerm)
	}
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
