package telegram

import (
	api "gopkg.in/telegram-bot-api.v4"
)

type Menu []MenuRow
type MenuRow []MenuItem
type MenuItem struct {
	api.KeyboardButton
	ReplyText string
	ReplyMenu []Row
}

type itemOption func(i *MenuItem)

func NewItem(btn api.KeyboardButton, options ...itemOption) MenuItem {
	i := MenuItem{KeyboardButton: btn}

	for _, opt := range options {
		opt(&i)
	}

	return i
}

func WithReplyText(text string) itemOption {
	return func(i *MenuItem) {
		i.ReplyText = text
	}
}

func WithReplyMenu(replyMenu ...Row) itemOption {
	return func(i *MenuItem) {
		i.ReplyMenu = replyMenu
	}
}

func FindItem(menu Menu, text string) *MenuItem {
	for _, row := range menu {
		for _, i := range row {
			if i.Text == text {
				return &i
			}
		}
	}
	return nil
}

func FindButton(menu Menu, data string) *Button {
	for _, row := range menu {
		for _, i := range row {
			if res := findButtonInRow(i.ReplyMenu, data); res != nil {
				return res
			}
		}
	}
	return nil
}

func findButtonInRow(queries []Row, data string) *Button {
	for _, qR := range queries {
		for _, q := range qR {
			if *q.CallbackData == data {
				return &q
			}

			for _, r := range q.Replies {
				if res := findButtonInRow(r.Menu, data); res != nil {
					return res
				}
			}
		}
	}
	return nil
}

func MenuToKeyboard(m ...MenuRow) api.ReplyKeyboardMarkup {
	butRows := make([][]api.KeyboardButton, 0, len(m))
	for _, row := range m {
		butRow := make([]api.KeyboardButton, 0, len(row))
		for _, b := range row {
			butRow = append(butRow, b.KeyboardButton)
		}
		butRows = append(butRows, butRow)
	}

	return api.NewReplyKeyboard(butRows...)
}

func MenuItemToInlineKeyboardMarkup(queryRows ...Row) api.InlineKeyboardMarkup {
	butRows := make([][]api.InlineKeyboardButton, 0, len(queryRows))
	for _, row := range queryRows {
		butRow := make([]api.InlineKeyboardButton, 0, len(row))
		for _, b := range row {
			butRow = append(butRow, b.InlineKeyboardButton)
		}
		butRows = append(butRows, butRow)
	}

	return api.NewInlineKeyboardMarkup(butRows...)
}
