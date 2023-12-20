package common

import (
	"check-price/src/core/constant"
	"context"
	"encoding/hex"
	"math/rand"
	"time"
)

// Detach returns a context that keeps all the values of its parent context
// but detaches from the cancellation and error handling.
func Detach(ctx context.Context) context.Context { return detachedContext{ctx} }

type detachedContext struct {
	parent context.Context
}

func (v detachedContext) Deadline() (time.Time, bool)       { return time.Time{}, false }
func (v detachedContext) Done() <-chan struct{}             { return nil }
func (v detachedContext) Err() error                        { return nil }
func (v detachedContext) Value(key interface{}) interface{} { return v.parent.Value(key) }

func CreateNewCtx() context.Context {
	ctx := context.TODO()
	var traceId string
	traceIdByte := make([]byte, 16)
	rand.Read(traceIdByte)
	traceId = hex.EncodeToString(traceIdByte[:])
	return context.WithValue(ctx, constant.TraceIdName, traceId)
}
