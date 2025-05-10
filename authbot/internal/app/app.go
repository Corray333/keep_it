package app

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"log/slog"
	"os"

	"github.com/corray333/keep_it_authbot/internal/storage"
	"github.com/corray333/keep_it_authbot/internal/types"
	"github.com/corray333/keep_it_authbot/internal/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Storage interface {
	SetUserRequest(query *types.CodeQuery) error
}

type App struct {
	Storage Storage
}

const (
	LoginRequestSignUp = iota + 1
	LoginRequestLogIn
	LoginRequestChangePassword
)

func New() *App {
	return &App{Storage: storage.New()}
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
				args_b64 := update.Message.CommandArguments()

				if args_b64 == "" {
					slog.Error("no start args error")
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, we have not enough info to verify your account. Did you get here through the app?🤨")
					bot.Send(msg)
					continue
				}

				decodedArgs, err := base64.StdEncoding.DecodeString(args_b64)
				if err != nil {
					slog.Error("decode error:", "error", err)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, we have some internal problems😢 Please, try to log in later.")
					bot.Send(msg)
					continue
				}
				var query types.CodeQuery
				if err := json.Unmarshal(decodedArgs, &query); err != nil {
					slog.Error("decoding error: ", "error", err)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, we have some internal problems😢 Please, try to log in later.")
					bot.Send(msg)
					continue
				}

				if query.Username == "" || query.Type == 0 {
					slog.Error("not full start args provided")
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, we have not enough info to verify your account. Did you get here through the app?🤨")
					bot.Send(msg)
					continue
				}

				query.TelegramID = update.FromChat().ID
				query.Code = utils.GenerateVerificationCode()

				if err := app.Storage.SetUserRequest(&query); err != nil {
					slog.Error("error while saving user request in redis: ", "error", err)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, we have some internal problems😢 Please, try to log in later.")
					bot.Send(msg)
					continue
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Your code is: "+query.Code)
				bot.Send(msg)
				continue
			}
		} else if update.Message.ForwardFrom != nil || update.Message.ForwardFromChat != nil {
			// link := ""
			// if update.Message.ForwardFromChat.UserName != "" {
			// 	// For public groups/channels
			// 	link = fmt.Sprintf("https://t.me/%s/%d", update.Message.ForwardFromChat.UserName, update.Message.ForwardFromMessageID)
			// } else {
			// 	// For private groups/channels
			// 	link = fmt.Sprintf("https://t.me/c/%d/%d", update.Message.ForwardFromChat.ID, update.Message.ForwardFromMessageID)
			// }

			// orig := types.Original{
			// 	Text: update.Message.ForwardFromChat.UserName,
			// 	Link: link,
			// }
			// marshalled, err := json.Marshal(orig)
			// if err != nil {
			// 	slog.Error("error while marshaling original message: ", "error", err)
			// 	continue
			// }
			// note := types.Note{
			// 	Original: string(marshalled),
			// 	Source:   "tg",
			// 	Type:     1,
			// 	Content:  update.Message.Text,
			// }
			// t, err := json.Marshal(note)

			// fmt.Println(err)
			// fmt.Println()
			// fmt.Println("Note: ", t)
			// fmt.Println()
		}
	}
}
