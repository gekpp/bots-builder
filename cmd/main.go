package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gekpp/bots-builder/internal/infra"
	"github.com/gekpp/bots-builder/questionnaire"
	"github.com/gekpp/bots-builder/users"
	"github.com/gekpp/env"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

var (
	argQuestionnaireID = env.MustString("QNR_ID")
	botToken           = env.MustString("TOKEN")
	debug              = env.GetBool("DEBUG", false)

	argDBHost     = env.MustString("DATABASE_HOST")
	argDBPort     = env.MustInt("DATABASE_PORT")
	argDBName     = env.MustString("DATABASE_NAME")
	argDBUsername = env.MustString("DATABASE_USER")
	argDBPassword = env.MustString("DATABASE_PASS")
	argDBTimeout  = env.GetInt("DATABASE_CONNECT_TIMEOUT", 15)

	qnr         questionnaire.Service
	userService users.Service
)

func main() {
	db := infra.MustConnectDB(argDBHost, argDBPort, argDBName, argDBUsername, argDBPassword, argDBTimeout, "disable")
	dbx := sqlx.NewDb(db, "postgres").Unsafe()
	qnr := questionnaire.New(dbx, uuid.MustParse(argQuestionnaireID))
	userService := users.New(dbx)
	_ = userService

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
		handleUpdate(ctx, qnr, userService, bot, update)
	}
}

func handleUpdate(ctx context.Context, qnr questionnaire.Service, usr users.Service, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	if update.Message == nil {
		return
	}

	u, err := usr.CreateOrGetTelegramUser(ctx, users.User{
		TelegramID:       fmt.Sprintf("%v", chatID),
		FirstName:        update.Message.Chat.FirstName,
		LastName:         update.Message.Chat.LastName,
		TelegramUserName: update.Message.Chat.UserName,
	})
	if err != nil {
		logrus.Error(err)
		return
	}

	if update.Message.IsCommand() {
		handleCommand(ctx, update.Message.Command(), u.ID, chatID, qnr, bot)
	} else {
		handleTextMessage(ctx, update.Message.Text, u.ID, chatID, qnr, bot)
	}
}

func handleCommand(
	ctx context.Context,
	command string,
	userID uuid.UUID,
	chatID int64,
	qnr questionnaire.Service,
	bot *tgbotapi.BotAPI) {

	var keyboard tgbotapi.ReplyKeyboardMarkup
	switch command {
	case "start":
		resp, err := qnr.Start(ctx, userID)
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
}

func handleTextMessage(
	ctx context.Context,
	text string,
	userID uuid.UUID,
	chatID int64,
	qnr questionnaire.Service,
	bot *tgbotapi.BotAPI) {
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
