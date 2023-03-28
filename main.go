package main

import (
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/sunchiii/telebot_gpt/client_api"
)

type Bot struct {
	Token string
}

func (b *Bot) Start() {
	topic := client_api.Topic{
		Model:       "gpt-3.5-turbo-0301",
		Temperature: 0.7,
	}
	bot, err := tgbotapi.NewBotAPI(b.Token)
	if err != nil {
		log.Fatal("can't connect telegram: ", err)
	}

	bot.Debug = true

	//fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		chatID := update.Message.Chat.ID
		actionTyping := tgbotapi.NewChatAction(chatID, tgbotapi.ChatTyping)
		if update.Message != nil {
			//send status typing
			bot.Send(actionTyping)
			//define msg
			Message := client_api.Message{
				Role:    "user",
				Content: update.Message.Text,
			}

			//append old msg and new msg
			topic.Messages = append(topic.Messages, Message)

			var choices = []client_api.Choice{}
			choices = client_api.Ask(topic)

			for _, message := range choices {
				newMsg := client_api.Message{}
				newMsg = message.Message
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, newMsg.Content)
				msg.ReplyToMessageID = update.Message.MessageID

				bot.Send(msg)
			}

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

	bot1.Start()

	time.Sleep(time.Hour)
}
