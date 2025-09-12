package nats

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Rasikrr/core/interfaces"
	"github.com/Rasikrr/core/log"

	"github.com/nats-io/nats.go"
)

type SubscriberOption func(*subscriber) error

type Subscriber interface {
	Subscribe(ctx context.Context, subject string, handler SubscriberHandler) error
	WithHandlers(handlers ...SubscriberHandler)
	interfaces.Closer
	interfaces.Starter
}

type SubscriberHandler interface {
	Handle(m *nats.Msg) error
	Subject() string
}

type subscriber struct {
	nc       *nats.Conn
	queue    string
	handlers []SubscriberHandler
}

func WithQueue(queue string) SubscriberOption {
	return func(s *subscriber) error {
		if s == nil {
			return errors.New("subscriber cannot be nil")
		}
		s.queue = queue
		return nil
	}
}

func NewSubscriber(addr string, options ...SubscriberOption) (Subscriber, error) {
	initNATSMetrics()
	nc, err := nats.Connect(addr)
	if err != nil {
		return nil, fmt.Errorf("connect to Nats %s error: %w", addr, err)
	}

	s := &subscriber{nc: nc}
	for _, opt := range options {
		if err = opt(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *subscriber) WithHandlers(handlers ...SubscriberHandler) {
	s.handlers = append(s.handlers, handlers...)
}

func (s *subscriber) Subscribe(ctx context.Context, subject string, handler SubscriberHandler) error {
	l := log.With(log.String("subject", subject))
	sub, err := s.nc.SubscribeSync(subject)
	if err != nil {
		return err
	}
	go func() {
		defer func() {
			if err := sub.Unsubscribe(); err != nil {
				log.Errorf(ctx, "unsubscribe error: %v", err)
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := sub.NextMsgWithContext(ctx)
				if err != nil {
					if errors.Is(err, context.Canceled) {
						return
					}
					log.Error(ctx, "next msg error", log.Err(err))
					continue
				}

				metrics.recvTotal.WithLabelValues(subject).Inc()
				if msg != nil && msg.Data != nil {
					metrics.recvBytes.WithLabelValues(subject).Observe(float64(len(msg.Data)))
				}

				start := time.Now()
				metrics.inflightReq.WithLabelValues(subject).Inc()
				if err := handler.Handle(msg); err != nil {
					l.Error(ctx, "handle message error", log.Err(err))
				}
				metrics.handlerSeconds.WithLabelValues(subject).Observe(time.Since(start).Seconds())
				metrics.inflightReq.WithLabelValues(subject).Dec()
			}
		}
	}()

	return nil
}

func (s *subscriber) SubscribeQueue(ctx context.Context, subject string, queue string, handler SubscriberHandler) error {
	l := log.With(log.String("subject", subject), log.String("queue", queue))

	_, err := s.nc.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		metrics.recvTotal.WithLabelValues(subject).Inc()
		if msg != nil && msg.Data != nil {
			metrics.recvBytes.WithLabelValues(subject).Observe(float64(len(msg.Data)))
		}

		start := time.Now()
		metrics.inflightReq.WithLabelValues(subject).Inc()
		if err := handler.Handle(msg); err != nil {
			l.Error(ctx, "handle message error", log.Err(err))
		}
		metrics.handlerSeconds.WithLabelValues(subject).Observe(time.Since(start).Seconds())
		metrics.inflightReq.WithLabelValues(subject).Dec()
	})
	log.Debugf(ctx, "subscribed to subject: %s, queue: %s\n", subject, queue)
	return err
}

func (s *subscriber) Start(ctx context.Context) error {
	for _, handler := range s.handlers {
		if err := s.SubscribeQueue(ctx, handler.Subject(), s.queue, handler); err != nil {
			return err
		}
	}
	return nil
}

func (s *subscriber) Close(ctx context.Context) error {
	s.nc.Close()
	log.Info(ctx, "nats subscriber closed")
	return nil
}
