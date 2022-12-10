package main

import (
	"context"
	"log"

	"github.com/gekpp/bots-builder/questionnaire"
	"github.com/gekpp/env"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

var (
	argQuestionnaireID = env.MustString("QNR_ID")
	botToken           = env.MustString("TOKEN")
	qnrService         = questionnaire.NewDummy(uuid.MustParse(argQuestionnaireID))
)

func main() {
	qnr := questionnaire.NewDummy(uuid.New())
	ctx := context.Background()
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = false
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		chatID := update.Message.Chat.ID
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			var resp questionnaire.StartResponse
			var keyboard tgbotapi.ReplyKeyboardMarkup
			switch update.Message.Command() {
			case "start":
				resp, err = qnr.Start(ctx, uuid.New())
				if err != nil {
					log.Fatal(err)
				}
			}

			err := sendMessage(bot, chatID, string(resp.Welcome), keyboard)
			if err != nil {
				log.Fatal(err)
			}

			keyboard = generateKeyboard(resp.Question)
			err = sendMessage(bot, chatID, string(resp.Question.Text), keyboard)
			if err != nil {
				log.Fatal(err)
			}

			continue
		}

		text := update.Message.Text

		resp, err := qnr.Answer(ctx, uuid.New(), questionnaire.Answer(text))
		if err != nil {
			log.Fatal(err)
		}

		keyboard := generateKeyboard(resp.Question)
		err = sendMessage(bot, chatID, string(resp.Question.Text), keyboard)
		if err != nil {
			log.Fatal(err)
		}
	}
}
