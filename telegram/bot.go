package telegram

import (
	"context"
	"github.com/Rasikrr/core/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/sync/semaphore"
)

type Bot interface {
	Start(ctx context.Context) error
	Close(ctx context.Context) error
	SendMessage(chatID int64, text string) error
	SendKeyboard(chatID int64, text string, buttons KeyBoard) error
	SetMessageHandler(handler func(update tgbotapi.Update))
}

type bot struct {
	sem     *semaphore.Weighted
	client  *tgbotapi.BotAPI
	handler func(update tgbotapi.Update)
}

func NewBot(token string, maxUserConcurrency int) (Bot, error) {
	client, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()

	var sem *semaphore.Weighted
	if maxUserConcurrency > 0 {
		log.Warnf(ctx, "Telegram bot initialized with max user concurrency: %d", maxUserConcurrency)
		sem = semaphore.NewWeighted(int64(maxUserConcurrency))
	} else {
		log.Warnf(ctx, "Telegram bot initialized with unlimited user concurrency")
	}

	return &bot{
		sem:    sem,
		client: client,
	}, nil
}

func (b *bot) SetMessageHandler(handler func(update tgbotapi.Update)) {
	b.handler = handler
}

func (b *bot) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := b.client.Send(msg)
	return err
}

func (b *bot) SendKeyboard(chatID int64, text string, buttons KeyBoard) error {
	keyboard := make([][]tgbotapi.KeyboardButton, len(buttons))
	for i, row := range buttons {
		btnRow := make([]tgbotapi.KeyboardButton, len(row))
		for j, label := range row {
			btnRow[j] = tgbotapi.NewKeyboardButton(label)
		}
		keyboard[i] = btnRow
	}
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(keyboard...)
	_, err := b.client.Send(msg)
	return err
}

// nolint: gocognit
func (b *bot) Start(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.client.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			return nil
		case update, ok := <-updates:
			if !ok {
				return nil
			}
			if update.Message != nil {
				if b.handler == nil {
					log.Warn(ctx,
						"No handler for message",
						log.Any("chat_id", update.Message.Chat.ID),
						log.String("bot_name", b.client.Self.UserName),
					)
					continue
				}
				go func() {
					if b.sem != nil {
						select {
						case <-ctx.Done():
							return
						default:
							err := b.sem.Acquire(ctx, 1)
							if err != nil {
								log.Error(ctx, "Failed to acquire semaphore", log.Err(err))
								return
							}
							defer b.sem.Release(1)
						}
					}
					b.handler(update)
				}()
			}
		}
	}
}

func (b *bot) Close(ctx context.Context) error {
	b.client.StopReceivingUpdates()
	log.Info(ctx, "Telegram bot closed", log.String("bot_name", b.client.Self.UserName))
	return nil
}
