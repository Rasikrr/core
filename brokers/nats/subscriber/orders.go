package subscriber

import (
	"context"

	"github.com/Rasikrr/core/log"
	"github.com/nats-io/nats.go"
)

type OrdersHandler struct {
	subject string
}

func NewOrdersHandler(subject string) *OrdersHandler {
	return &OrdersHandler{
		subject: subject,
	}
}

func (h *OrdersHandler) Subject() string {
	return h.subject
}

func (h *OrdersHandler) Handle(m *nats.Msg) error {
	ctx := context.Background()
	log.Infof(ctx, "handled in NATS: %+v", m)
	return nil
}
