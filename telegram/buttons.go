package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewMenu(row ...Row) []Row {
	return row
}

func NewSingleButtonMenu(btn Button) []Row {
	return []Row{NewRow(btn)}
}

type Row []Button

func NewRow(btn ...Button) Row {
	return btn
}

type Button struct {
	tgbotapi.InlineKeyboardButton
	Replies []Reply
}

type Reply struct {
	Menu     []Row //todo Сделать структуру Text+[]QueryRows
	Text     string
	Location *tgbotapi.Location
	Survey   *SurveyConfig
}

type buttonOption func(b *Button)

func WithSurvey(s SurveyConfig) buttonOption {
	return func(b *Button) {
		b.Replies = append(b.Replies, Reply{Survey: &s})
	}
}

func WithTextReply(text string) buttonOption {
	return func(b *Button) {
		b.Replies = append(b.Replies, Reply{Text: text})
	}
}

func WithMenuReply(text string, menu []Row) buttonOption {
	return func(b *Button) {
		b.Replies = append(b.Replies, Reply{Text: text, Menu: menu})
	}
}

func NewButton(text, data string, options ...buttonOption) Button {
	res := Button{
		InlineKeyboardButton: tgbotapi.NewInlineKeyboardButtonData(text, data),
	}

	for _, opt := range options {
		opt(&res)
	}

	return res
}
