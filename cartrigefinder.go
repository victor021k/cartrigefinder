package main

import (
	"cartrigefinder/db"
	"io/ioutil"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func getToken() string {
	data, err := ioutil.ReadFile(".token")
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSpace(string(data))
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	bot, err := tgbotapi.NewBotAPI(getToken())
	if err != nil {
		log.Fatal(err)
	}

	// bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	log.Println(db.Customers)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		//log.Printf("%+v\n", update.Message.Chat)
		if update.Message.Text == "/start" {
			log.Println("Человек ввёл start")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			//msg.ReplyToMessageID = update.Message.MessageID
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Заменить картридж")))
		}

		bot.Send(msg)
	}
}
