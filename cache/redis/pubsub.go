package redis

import (
	"context"

	"github.com/Rasikrr/core/log"
	goredis "github.com/redis/go-redis/v9"
)

type Subscription struct {
	pubSub *PubSub
}

func (s *Subscription) Close() error {
	return s.pubSub.Close()
}

func (c *Client) Publish(ctx context.Context, channel string, data any) error {
	return c.client.Publish(ctx, channel, data).Err()
}

// Subscribe subscribes to channels and returns PubSub for manual control.
// You must call Close() when done to release resources.
//
// Example:
//
//	pubsub := cache.Subscribe(ctx, "notifications", "alerts")
//	defer pubsub.Close()
//
//	for msg := range pubsub.Channel() {
//	    fmt.Println("Received:", msg.Payload)
//	}
func (c *Client) Subscribe(ctx context.Context, channels ...string) *PubSub {
	return c.client.Subscribe(ctx, channels...)
}

func (c *Client) SubscribeWithHandler(
	ctx context.Context,
	handler func(msg *Message) error,
	channels ...string,
) (*Subscription, error) {
	pubSub := c.client.Subscribe(ctx, channels...)

	if _, err := pubSub.Receive(ctx); err != nil {
		pubSub.Close()
		return nil, err
	}

	go func() {
		defer pubSub.Close()
		ch := pubSub.Channel()

		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				if err := handler(msg); err != nil {
					c.logger.Error(ctx, "pubsub handler error", log.Err(err))
				}
			}
		}
	}()

	return &Subscription{
		pubSub: pubSub,
	}, nil
}

func (c *Client) PSubscribe(ctx context.Context, patterns ...string) *goredis.PubSub {
	return c.client.PSubscribe(ctx, patterns...)
}
