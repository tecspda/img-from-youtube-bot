package models

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func GetKb(chat_id int64, db *Database) tgbotapi.ReplyKeyboardMarkup {
	buttons := []tgbotapi.KeyboardButton{
		tgbotapi.NewKeyboardButton("Помощь"),
		tgbotapi.NewKeyboardButton("Наши проекты"),
	}

	// Разделение кнопок на два ряда, если более двух
	rows := make([][]tgbotapi.KeyboardButton, 0)
	row := make([]tgbotapi.KeyboardButton, 0)

	for _, button := range buttons {
		// Добавить кнопку в текущий ряд
		row = append(row, button)
		// Если ряд заполнен, добавить его в общий массив и начать новый ряд
		if len(row) == 2 {
			rows = append(rows, row)
			row = make([]tgbotapi.KeyboardButton, 0)
		}
	}

	// Если остались кнопки, добавить последний ряд
	if len(row) > 0 {
		rows = append(rows, row)
	}

	keyboard := tgbotapi.NewReplyKeyboard(rows...)

	return keyboard
}
