package pubsub

type Subscriber interface {
	Handle(event Event)
}
