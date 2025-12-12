package nats

import (
	"context"

	"github.com/Rasikrr/core/sentry"
)

func setSentryHubAndScope(ctx context.Context, msg *Msg, queue string) context.Context {
	if !sentry.Enabled() {
		return ctx
	}

	ctx = sentry.SetHubOnCtx(ctx)
	hub := sentry.GetHubFromContext(ctx)

	// Устанавливаем контекст с информацией о NATS
	natsContext := map[string]interface{}{
		"type": "subscriber",
	}

	if queue != "" {
		natsContext["queue"] = queue
	}

	if msg != nil {
		natsContext["subject"] = msg.Subject

		if msg.Data != nil {
			natsContext["message_size"] = len(msg.Data)
		}
		if msg.Reply != "" {
			natsContext["reply_to"] = msg.Reply
		}
		if msg.Header != nil {
			natsContext["headers"] = msg.Header
		}
	}

	hub.Scope().SetContext("nats", natsContext)

	// Добавляем теги для быстрого поиска в Sentry
	hub.Scope().SetTag("transport", "nats")
	if msg != nil {
		hub.Scope().SetTag("nats.subject", msg.Subject)
	}
	if queue != "" {
		hub.Scope().SetTag("nats.queue", queue)
	}

	return ctx
}
