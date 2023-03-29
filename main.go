package main

import (
	"fmt"
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

	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	var chatHistory []int

	for update := range updates {
		chatID := update.Message.Chat.ID
		actionTyping := tgbotapi.NewChatAction(chatID, tgbotapi.ChatTyping)

		if update.Message != nil {
			//keep chat id
			chatHistory = append(chatHistory, update.Message.MessageID)
			//send status typing
			bot.Send(actionTyping)
			//check command
			if update.Message.IsCommand() {
				if update.Message.Text == "/start" {
					chatHistory = []int{}
					topic.Messages = []client_api.Message{}
					continue
				}
				if update.Message.Text == "/clearchat" {
					//clear chatHistory
					for _, msgId := range chatHistory {
						deleteMessage := tgbotapi.NewDeleteMessage(chatID, msgId)
						bot.Send(deleteMessage)
					}

					continue
				}
				if update.Message.Text == "/newtopic" {
					//clear topice
					topic.Messages = []client_api.Message{}
					//clear chat
					for _, msgId := range chatHistory {
						deleteMessage := tgbotapi.NewDeleteMessage(chatID, msgId)
						bot.Send(deleteMessage)
					}
					continue
				}
			}
			//define msg
			Message := client_api.Message{
				Role:    "user",
				Content: update.Message.Text,
			}

			//append old msg and new msg
			topic.Messages = append(topic.Messages, Message)

			//get choices from api
			var choices = []client_api.Choice{}
			choices = client_api.Ask(topic)

			//reply each message
			for _, message := range choices {
				newMsg := client_api.Message{}
				newMsg = message.Message
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, newMsg.Content)
				msg.ReplyToMessageID = update.Message.MessageID

				mId, _ := bot.Send(msg)
				chatHistory = append(chatHistory, mId.MessageID)
			}

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
	sunday := Bot{Token: os.Getenv("TG_SUNDAY")}
	humyai := Bot{Token: os.Getenv("TG_HUMYAI")}
	aiChatBot := Bot{Token: os.Getenv("TG_AICHATBOT")}

	go sunday.Start()
	go humyai.Start()
	go aiChatBot.Start()

	time.Sleep(time.Hour * 8000)
}
