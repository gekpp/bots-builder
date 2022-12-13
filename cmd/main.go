package main

import (
	"context"
	"log"

	"github.com/gekpp/bots-builder/questionnaire"
	"github.com/gekpp/env"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	argQuestionnaireID = env.MustString("QNR_ID")
	botToken           = env.MustString("TOKEN")
	debug              = env.GetBool("DEBUG", false)
	qnrService         = questionnaire.NewDummy(uuid.MustParse(argQuestionnaireID))
)

func handleMessage(qnr questionnaire.Service, ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	if update.Message == nil {
		return
	}

	if update.Message.IsCommand() {
		var keyboard tgbotapi.ReplyKeyboardMarkup
		switch update.Message.Command() {
		case "start":
			resp, err := qnr.Start(ctx, uuid.New())
			if err != nil {
				logrus.Error(err)
			}
			err = sendMessagePlain(bot, chatID, string(resp.Welcome))
			if err != nil {
				logrus.Error(err)
			}

			keyboard = generateKeyboard(resp.Question)
			err = sendMessageWithKeyboard(bot, chatID, string(resp.Question.Text), keyboard)
			if err != nil {
				logrus.Error(err)
			}
		default:
			err := sendMessagePlain(bot, chatID, "Такой команды нет.")
			if err != nil {
				logrus.Error(err)
			}
		}

		return
	}

	text := update.Message.Text

	resp, err := qnr.Answer(ctx, uuid.New(), questionnaire.Answer(text))
	if err != nil {
		logrus.Error(err)
	}

	keyboard := generateKeyboard(resp.Question)
	err = sendMessageWithKeyboard(bot, chatID, string(resp.Question.Text), keyboard)
	if err != nil {
		logrus.Error(err)
	}

	return
}

func main() {
	qnr := questionnaire.NewDummy(uuid.New())
	ctx := context.Background()
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = debug
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		handleMessage(qnr, ctx, bot, update)
	}
}
