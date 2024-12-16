package repo

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	apiToken = "api token"
	chatID   = 1234256
)

type TelegramBot struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

func NewTelegramBot() (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к Telegram API: %w", err)
	}

	return &TelegramBot{
		bot:    bot,
		chatID: chatID,
	}, nil
}

func (t *TelegramBot) SendMessage(message string) error {
	msg := tgbotapi.NewMessage(t.chatID, message)
	_, err := t.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("ошибка отправки сообщения: %w", err)
	}
	return nil
}
