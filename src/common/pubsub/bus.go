package pubsub

import (
	"check-price/src/common/log"
)

type EventBus struct {
	producer     Producer
	eventMapping map[string][]Subscriber
}

func NewEventBus(publisher *publisher) *EventBus {
	return &EventBus{
		producer:     publisher,
		eventMapping: make(map[string][]Subscriber),
	}
}

func (b *EventBus) Subscribe(event Event, subscriber ...Subscriber) {
	if _, exist := b.eventMapping[event.GetName()]; exist == false {
		b.eventMapping[event.GetName()] = make([]Subscriber, 0)
	}
	b.eventMapping[event.GetName()] = append(b.eventMapping[event.GetName()], subscriber...)
}

func (b *EventBus) Run() {
	for {
		event := <-b.producer.ProduceEvent()
		log.Info(event.GetContext(), "event [%s] was fired", event.GetName())
		for _, subscriber := range b.eventMapping[event.GetName()] {
			go subscriber.Handle(event)
		}
	}
}
