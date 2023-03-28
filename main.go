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
	topic := Topic{
		Model:       "gpt-3.5-turbo-0301",
		Temperature: "0.7",
		Messages:    []Message{},
	}
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
			//send status typing
			bot.Send(actionTyping)
			//define msg
			Message := Message{
				Role:    "user",
				Content: update.Message.Text,
			}

			topic.Messages = append(topic.Messages, Message)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
			continue
		}
		if update.Message.IsCommand() {
			if update.Message.Text == "/clearChat" {

				continue
			}
			if update.Message.Text == "/newTopic" {

				continue
			}
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
