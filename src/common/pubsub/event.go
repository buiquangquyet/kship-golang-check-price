package pubsub

import "context"

type Event interface {
	GetName() string
	Payload() interface{}
	GetContext() context.Context
}
