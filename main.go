package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

type Bot struct {
	Token string
}

func (b *Bot) Start() {
	bot, err := tgbotapi.NewBotAPI(b.Token)
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		chatID := update.Message.Chat.ID
		actionTyping := tgbotapi.NewChatAction(chatID, tgbotapi.ChatTyping)
		if update.Message != nil {
			bot.Send(actionTyping)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
			continue
		}

	}
}

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("can't found .env file")
	}
}

func main() {
	bot1 := Bot{Token: os.Getenv("TG_SUNDAY")}

	go bot1.Start()

	time.Sleep(8760 * time.Hour)
}
