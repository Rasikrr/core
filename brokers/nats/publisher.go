package nats

import (
	"context"
	"time"

	"github.com/Rasikrr/core/interfaces"
	"github.com/Rasikrr/core/log"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type Publisher interface {
	Publish(ctx context.Context, subject string, m proto.Message) error
	interfaces.Closer
}

type publisher struct {
	conn *nats.Conn
}

func NewPublisher(addr string) (Publisher, error) {
	initNATSMetrics()
	conn, err := nats.Connect(
		addr,
		nats.MaxReconnects(-1), // бесконечные реконнекты
		nats.ReconnectWait(time.Second),
	)
	if err != nil {
		return nil, err
	}
	return &publisher{
		conn: conn,
	}, nil
}

func (p *publisher) Publish(_ context.Context, subject string, m proto.Message) error {
	bb, err := proto.Marshal(m)
	if err != nil {
		return err
	}

	if len(bb) > 0 {
		metrics.pubBytes.WithLabelValues(subject).Observe(float64(len(bb)))
	}
	metrics.pubTotal.WithLabelValues(subject).Inc()

	msg := &nats.Msg{
		Subject: subject,
		Data:    bb,
		Header:  make(nats.Header),
	}
	msg.Header.Set("Content-Type", "application/protobuf")
	msg.Header.Set("Content-Encoding", "binary")

	return p.conn.PublishMsg(msg)
}

func (p *publisher) Close(ctx context.Context) error {
	p.conn.Close()
	log.Info(ctx, "nats publisher closed")
	return nil
}
