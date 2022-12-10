package main

import (
	"reflect"

	"github.com/gekpp/bots-builder/questionnaire"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func generateKeyboard(qnrResponse questionnaire.Question) tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard()

	for _, item := range qnrResponse.AnswerOptions {
		btnslc := []tgbotapi.KeyboardButton{}
		btnslc = append(btnslc, tgbotapi.NewKeyboardButton(string(item)))
		keyboard.Keyboard = append(keyboard.Keyboard, btnslc)
		keyboard.OneTimeKeyboard = true
	}

	return keyboard
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, msgText string, keyboard tgbotapi.ReplyKeyboardMarkup) error {
	msg := tgbotapi.NewMessage(chatID, msgText)
	if !reflect.DeepEqual(keyboard, tgbotapi.ReplyKeyboardMarkup{}) {
		msg.ReplyMarkup = keyboard
	}

	if _, err := bot.Send(msg); err != nil {
		return err
	}

	return nil
}
