package log

import (
	"check-price/src/core/constant"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

const callerSkip = 2

type logger struct {
	zap *zap.SugaredLogger
}

func SyslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func CustomLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

func NewLogger() {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:   "message",
		LevelKey:     "level",
		TimeKey:      "time",
		CallerKey:    "caller",
		EncodeCaller: zapcore.FullCallerEncoder,
		EncodeTime:   SyslogTimeEncoder,
		EncodeLevel:  CustomLevelEncoder,
	}

	var encoder zapcore.Encoder
	var level zapcore.Level
	if constant.IsProdEnv() {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
		level = zap.InfoLevel
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
		level = zap.DebugLevel
	}
	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stderr), level)
	//set log
	globalLogger = &logger{
		zap: zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel), zap.AddCallerSkip(callerSkip)).Sugar(),
	}
	return
}

func (l *logger) Info(msg string, args ...interface{}) {
	l.zap.Infof(msg, args...)
}

func (l *logger) Debug(msg string, args ...interface{}) {
	l.zap.Debugf(msg, args...)
}

func (l *logger) Warn(msg string, args ...interface{}) {
	l.zap.Warnf(msg, args...)
}

func (l *logger) Error(msg string, args ...interface{}) {
	l.zap.Errorf(msg, args...)
}

func (l *logger) Fatal(msg string, args ...interface{}) {
	l.zap.Fatalf(msg, args...)
}

func (l *logger) GetZap() *zap.SugaredLogger {
	return l.zap
}
