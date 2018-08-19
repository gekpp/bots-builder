package telegram

import (
	"gopkg.in/telegram-bot-api.v4"
	"fmt"
)

type Sender struct {
	bot *tgbotapi.BotAPI
}

func NewSender(bot *tgbotapi.BotAPI) *Sender {
	s := Sender{
		bot: bot,
	}
	return &s
}

type messageOption func(m *tgbotapi.BaseChat)

func removeKeyboard(m *tgbotapi.BaseChat) {
	m.ReplyMarkup = tgbotapi.NewRemoveKeyboard(false)
}

func MarkDown(m *tgbotapi.MessageConfig) {
	m.ParseMode = "Markdown"
}

func forceReply(m *tgbotapi.BaseChat) {
	m.ReplyMarkup = tgbotapi.ForceReply{}
}

func InlineKeyboard(k tgbotapi.BaseChat) messageOption {
	return func(m *tgbotapi.BaseChat) {
		m.ReplyMarkup = k
	}
}

func requestContact(yesButton string, noButton string) func(m *tgbotapi.MessageConfig) {
	return func(m *tgbotapi.MessageConfig) {
		keyboard := MenuToKeyboard(
			MenuRow{
				NewItem(tgbotapi.NewKeyboardButtonContact(yesButton)),
				NewItem(tgbotapi.NewKeyboardButton(noButton)),
			},
		)
		m.BaseChat.ReplyMarkup = keyboard
	}
}

func text(update tgbotapi.Update) string {
	if update.CallbackQuery != nil {
		return update.CallbackQuery.Data
	}

	if update.Message != nil {
		res := update.Message.Text
		if update.Message.Contact != nil {
			res += update.Message.Contact.PhoneNumber
		}
		return res
	}
	return ""
}

func (s *Sender) SendTextMessage(chatId int64, msg string, options ... messageOption) error {
	message := tgbotapi.NewMessage(chatId, msg)
	for _, opt := range options {
		opt(&message.BaseChat)
	}
	_, err := s.bot.Send(message)
	return err
}

func (s *Sender) SendPhoto(chatId int64, photoUrl string, options ...messageOption) error {
	message := tgbotapi.NewPhotoUpload(chatId, photoUrl)
	for _, opt := range options {
		opt(&message.BaseChat)
	}
	_, err := s.bot.Send(message)
	return err
}

func (s *Sender) SendMenuItemReply(item MenuItem, chatId int64) error {
	keyboard := MenuItemToInlineKeyboardMarkup(item.ReplyMenu...)

	reply := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:      chatId,
			ReplyMarkup: keyboard,
		},
		Text: item.ReplyText,
	}
	_, err := s.bot.Send(&reply)
	return fmt.Errorf("unable to send menu item: %v", err)
}

func (s *Sender) SendMenu(m Menu, chatId int64) error {
	keyboard := MenuToKeyboard(m...)

	reply := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:      chatId,
			ReplyMarkup: keyboard,
		},
		Text: `Чем я могу Вам помочь?
Для перехода в главное меню, нужно нажать кнопку "Меню" внизу экрана`,
	}

	_, err := s.bot.Send(&reply)
	return fmt.Errorf("unable to send menu: %v", err)
}
