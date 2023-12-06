package log

import (
	"check-price/src/common/configs"
	"check-price/src/core/constant"
	"check-price/src/core/enums"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
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
	cores := make([]zapcore.Core, 0)
	cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stderr), level))
	if configs.Get().Log.EnableFile {
		cores = append(cores, zapcore.NewCore(encoder, getWriteSyncer(), level))
	}
	tee := zapcore.NewTee(cores...)
	globalLogger = &logger{
		zap: zap.New(tee, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel), zap.AddCallerSkip(callerSkip)).Sugar(),
	}
	return
}

func (l *logger) Info(msg string, traceId enums.TraceId, merchant enums.Merchant) {
	l.zap.Infow(msg, zap.Inline(traceId), zap.Inline(merchant))
}

func (l *logger) Debug(msg string, traceId enums.TraceId, merchant enums.Merchant) {
	l.zap.Debugf(msg, zap.Inline(traceId), zap.Inline(merchant))
}

func (l *logger) Warn(msg string, traceId enums.TraceId, merchant enums.Merchant) {
	l.zap.Warnf(msg, zap.Inline(traceId), zap.Inline(merchant))
}

func (l *logger) Error(msg string, traceId enums.TraceId, merchant enums.Merchant) {
	l.zap.Errorf(msg, zap.Inline(traceId), zap.Inline(merchant))
}

func (l *logger) Fatal(msg string, args ...interface{}) {
	l.zap.Fatalf(msg, args...)
}

func (l *logger) GetZap() *zap.SugaredLogger {
	return l.zap
}

func getWriteSyncer() zapcore.WriteSyncer {
	cf := configs.Get().Log
	lumberJackLogger := &lumberjack.Logger{
		Filename:   cf.File,
		MaxSize:    cf.MaxSize,
		MaxAge:     cf.MaxAge,
		MaxBackups: cf.MaxBackups,
		Compress:   true,
	}
	return zapcore.AddSync(lumberJackLogger)

}
