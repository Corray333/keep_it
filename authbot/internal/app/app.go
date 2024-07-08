package app

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"log/slog"
	"os"

	"github.com/Corray333/authbot/internal/storage"
	"github.com/Corray333/authbot/internal/types"
	"github.com/Corray333/authbot/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Storage interface {
	SetUserRequest(query *types.CodeQuery) error
}

type App struct {
	Storage Storage
}

const (
	LoginRequestLogIn = iota
	LoginRequestSignUp
	LoginRequestChangePassword
)

func New() *App {
	return &App{Storage: storage.NewStorage()}
}

func (app *App) Run() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // игнорировать все не-сообщения
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				// Получение аргументов командной строки
				args_b64 := update.Message.CommandArguments()

				if args_b64 == "" {
					slog.Error("no start args error")
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, we have not enough info to verify your account. Did you get here through the app?🤨")
					bot.Send(msg)
					continue
				}

				decodedArgs, err := base64.StdEncoding.DecodeString(args_b64)
				if err != nil {
					slog.Error("decode error:" + err.Error())
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, we have some internal problems😢 Please, try to log in later.")
					bot.Send(msg)
					return
				}
				var query types.CodeQuery
				if err := json.Unmarshal(decodedArgs, &query); err != nil {
					slog.Error("decoding error: " + err.Error())
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, we have some internal problems😢 Please, try to log in later.")
					bot.Send(msg)
					return
				}

				if query.Username == "" || query.TypeID == 0 {
					slog.Error("not full start args provided")
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, we have not enough info to verify your account. Did you get here through the app?🤨")
					bot.Send(msg)
					continue
				}

				query.TG = update.Message.From.UserName
				query.Code = utils.GenerateVerificationCode()

				if err := app.Storage.SetUserRequest(&query); err != nil {
					slog.Error("error while saving user request in redis: " + err.Error())
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, we have some internal problems😢 Please, try to log in later.")
					bot.Send(msg)
					continue
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Your code is: "+query.Code)
				bot.Send(msg)
				continue

			}
		}
	}
}
