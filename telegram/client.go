package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
)

type Client interface {
	SendText(ctx context.Context, chatID int64, text string) error
	SendFile(ctx context.Context, chatID int64, file io.Reader, fileName, caption string) error
}

type client struct {
	bot *tgbotapi.BotAPI
}

func NewTelegramClient(token string) (Client, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &client{bot: bot}, nil
}

func (c *client) SendText(_ context.Context, chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := c.bot.Send(msg)
	return err
}

func (c *client) SendFile(_ context.Context, chatID int64, file io.Reader, fileName string, caption string) error {
	doc := tgbotapi.NewDocument(chatID, tgbotapi.FileReader{
		Name:   fileName,
		Reader: file,
	})
	doc.Caption = caption
	_, err := c.bot.Send(doc)
	return err
}
