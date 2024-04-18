package types

import (
	"context"
	"inttest-runtime/pkg/embedded"
	"inttest-runtime/pkg/utils"
	"log"
	"time"

	"golang.org/x/sync/errgroup"
)

type MockBroker struct {
	eventBus IMockBrokerPubSub
	rules    []iMockBrokerRule

	exec *PyPrecompiledExecutor
}

type IMockBrokerPubSub interface {
	Publish(ctx context.Context, topic string, message []byte) error
	Subscribe(ctx context.Context, topic string, consumer func(message []byte) error)
}

type iMockBrokerRule interface {
	GetMsg() ([]byte, error)
	GetCommon() brokerRuleCommon
}

type brokerRuleCommon struct {
	interval        time.Duration
	sendImmediately bool
	topic           string
}

func (c brokerRuleCommon) GetCommon() brokerRuleCommon {
	return c
}

type mockBrokerRule_Raw struct {
	brokerRuleCommon
	message []byte
}

func (r mockBrokerRule_Raw) GetMsg() ([]byte, error) {
	return r.message, nil
}

type mockBrokerRule_Programmable struct {
	brokerRuleCommon

	msgGenerator embedded.CodeSnippet
	executor     *PyPrecompiledExecutor
}

func (r mockBrokerRule_Programmable) GetMsg() ([]byte, error) {
	msg, err := r.executor.ExecFunc(r.msgGenerator)
	if err != nil {
		return nil, err
	}
	msgStr, err := msg.ToString()
	if err != nil {
		return nil, err
	}

	return utils.S2B(msgStr), nil
}

func NewMockBroker(funcExec *PyPrecompiledExecutor, eventBus IMockBrokerPubSub) *MockBroker {
	return &MockBroker{
		exec:     funcExec,
		eventBus: eventBus,
	}
}

func (b *MockBroker) AddStubRule(interval time.Duration, sendImmediately bool, topic string, msg []byte) {
	b.rules = append(b.rules, mockBrokerRule_Raw{
		brokerRuleCommon: brokerRuleCommon{
			interval:        interval,
			sendImmediately: sendImmediately,
			topic:           topic,
		},

		message: msg,
	})
}

func (b *MockBroker) AddProgrammableRule(interval time.Duration, sendImmediately bool, topic string, generator embedded.CodeSnippet) {
	b.rules = append(b.rules, mockBrokerRule_Programmable{
		brokerRuleCommon: brokerRuleCommon{
			interval:        interval,
			sendImmediately: sendImmediately,
			topic:           topic,
		},

		msgGenerator: generator,
		executor:     b.exec,
	})
}

func (b *MockBroker) Start(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, r := range b.rules {
		rule := r
		g.Go(func() error {
			if rule.GetCommon().sendImmediately {
				msg, err := rule.GetMsg()
				if err != nil {
					log.Printf("error sending msg to topic: %v\n", err)
				}
				if err := b.eventBus.Publish(ctx, rule.GetCommon().topic, msg); err != nil {
					log.Printf("error sending msg to topic: %v\n", err)
				}
			}
			for {
				time.Sleep(rule.GetCommon().interval)
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					msg, err := rule.GetMsg()
					if err != nil {
						log.Printf("error sending msg to topic: %v\n", err)
					}
					if err := b.eventBus.Publish(ctx, rule.GetCommon().topic, msg); err != nil {
						log.Printf("error sending msg to topic: %v\n", err)
					}
				}
			}
		})
	}
	return g.Wait()
}
