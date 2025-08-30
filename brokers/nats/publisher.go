package nats

import (
	"context"
	"time"

	"github.com/Rasikrr/core/interfaces"
	"github.com/Rasikrr/core/log"
	"github.com/Rasikrr/core/metrics"
	nats2 "github.com/Rasikrr/core/metrics/nats"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type Publisher interface {
	Publish(ctx context.Context, subject string, m proto.Message) error
	interfaces.Closer
}

type publisher struct {
	conn    *nats.Conn
	metrics nats2.PublisherMetrics
}

func NewPublisher(
	addr string,
	metricer metrics.Metricer,
) (Publisher, error) {
	conn, err := nats.Connect(addr)
	if err != nil {
		return nil, err
	}
	pub := &publisher{
		conn: conn,
	}

	pub.metrics = nats2.NewNATSPublisherMetrics(metricer)

	return pub, nil
}

func (p *publisher) Publish(_ context.Context, subject string, m proto.Message) error {
	start := time.Now()
	defer func() {
		p.metrics.ObserveDuration(subject, time.Since(start).Seconds())
	}()

	bb, err := proto.Marshal(m)
	if err != nil {
		p.metrics.IncError(subject)
		return err
	}

	msg := &nats.Msg{
		Subject: subject,
		Data:    bb,
		Header:  make(nats.Header),
	}
	msg.Header.Set("Content-Type", "application/protobuf")
	msg.Header.Set("Content-Encoding", "binary")

	err = p.conn.PublishMsg(msg)
	if err != nil {
		p.metrics.IncError(subject)
	}
	return err
}

func (p *publisher) Close(ctx context.Context) error {
	p.conn.Close()
	log.Info(ctx, "nats publisher closed")
	return nil
}
