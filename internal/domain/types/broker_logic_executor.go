package types

import "context"

type MockBroker struct {
	eventBus IMockBrokerPubSub
}

type IMockBrokerPubSub interface {
	Publish(ctx context.Context, message []byte) error
	Subscribe(ctx context.Context, topic string, consumer func(message []byte) error)
}

func NewMockBroker()
