package repo

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	apiToken = ""
	chatID   = 1234567
)

type TelegramBot struct {
	bot           *tgbotapi.BotAPI
	chatID        int64
	ventilationOn bool // Флаг состояния вентиляции
}

func NewTelegramBot() (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к Telegram API: %w", err)
	}

	return &TelegramBot{
		bot:           bot,
		chatID:        chatID,
		ventilationOn: false, // Вентиляция выключена по умолчанию
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

// Отправка сообщения с кнопками
func (t *TelegramBot) SendControlButtons() error {
	// Создание инлайн-кнопок
	buttonOn := tgbotapi.NewInlineKeyboardButtonData("Включить вентиляцию", "turn_on")
	buttonOff := tgbotapi.NewInlineKeyboardButtonData("Выключить вентиляцию", "turn_off")

	// Формирование клавиатуры
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(buttonOn, buttonOff),
	)

	// Отправка сообщения с клавиатурой
	msg := tgbotapi.NewMessage(t.chatID, "Управление вентиляцией:")
	msg.ReplyMarkup = keyboard

	_, err := t.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("ошибка отправки сообщения с кнопками: %w", err)
	}
	return nil
}

// Обработка обновлений (нажатий кнопок)
func (t *TelegramBot) HandleUpdates() (error, bool) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			// Обработка нажатий инлайн-кнопок
			data := update.CallbackQuery.Data
			var response string

			switch data {
			case "turn_on":
				t.ventilationOn = true
				response = "Вентиляция воздуха включена."
			case "turn_off":
				t.ventilationOn = false
				response = "Вентиляция воздуха отключена."
			default:
				response = "Неизвестная команда."
			}

			// Отправка ответа пользователю
			msg := tgbotapi.NewMessage(t.chatID, response)
			_, err := t.bot.Send(msg)
			if err != nil {
				return fmt.Errorf("ошибка отправки ответа: %w", err), false
			}

			// Уведомляем Telegram об обработке нажатия
			t.bot.Request(tgbotapi.NewCallback(update.CallbackQuery.ID, "Команда выполнена"))
		}
	}

	return nil, t.ventilationOn
}
