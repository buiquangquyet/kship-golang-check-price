package enums

import "go.uber.org/zap/zapcore"

type Merchant string

func (m Merchant) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("merchant", string(m))
	return nil
}
