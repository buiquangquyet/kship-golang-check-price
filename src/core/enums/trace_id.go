package enums

import "go.uber.org/zap/zapcore"

type TraceId string

func (t TraceId) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("trace_id", string(t))
	return nil
}
