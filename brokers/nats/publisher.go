package nats

import (
	"context"
	"github.com/Rasikrr/core/interfaces"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"log"
)

type Publisher interface {
	Publish(ctx context.Context, subject string, m proto.Message) error
	interfaces.Closer
}

type publisher struct {
	conn *nats.Conn
}

func NewPublisher(addr string) (Publisher, error) {
	conn, err := nats.Connect(addr)
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

	msg := &nats.Msg{
		Subject: subject,
		Data:    bb,
		Header:  make(nats.Header),
	}
	msg.Header.Set("Content-Type", "application/protobuf")
	msg.Header.Set("Content-Encoding", "binary")

	return p.conn.PublishMsg(msg)
}

func (p *publisher) Close(_ context.Context) error {
	p.conn.Close()
	log.Println("nats publisher closed")
	return nil
}
