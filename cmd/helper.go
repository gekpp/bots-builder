package main

import (
	"errors"

	"github.com/gekpp/bots-builder/questionnaire"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func generateKeyboard(options []questionnaire.Answer) tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard()

	for _, item := range options {
		btnslc := []tgbotapi.KeyboardButton{}
		btnslc = append(btnslc, tgbotapi.NewKeyboardButton(string(item)))
		keyboard.Keyboard = append(keyboard.Keyboard, btnslc)
	}
	keyboard.OneTimeKeyboard = true

	return keyboard
}

func sendMessageWithKeyboard(bot *tgbotapi.BotAPI, chatID int64, msgText string, keyboard tgbotapi.ReplyKeyboardMarkup) error {
	msg := tgbotapi.NewMessage(chatID, msgText)
	msg.ReplyMarkup = keyboard
	if len(keyboard.Keyboard) == 0 {
		return errors.New("keyboard shouldn't be empty")
	}

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}

func sendMessagePlain(bot *tgbotapi.BotAPI, chatID int64, msgText string) error {
	msg := tgbotapi.NewMessage(chatID, msgText)

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}
