package pubsub

var _publisher Publisher

func RegisterGlobal(publisher *publisher) {
	_publisher = publisher
}

func Publish(event Event) {
	_publisher.Publish(event)
}
