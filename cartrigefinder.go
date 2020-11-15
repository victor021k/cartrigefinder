package main

import (
	"cartrigefinder/db"
	"fmt"
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

func getCustomerNameById(id int64) string {
	name := ""
	for _, customer := range db.Customers {
		if customer.ID == id {
			name = customer.Name
			break
		}
	}

	return name
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	bot, err := tgbotapi.NewBotAPI(getToken())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

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
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Здравствуйте, я помогу заменить Вам картридж")
			//msg.ReplyToMessageID = update.Message.MessageID
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Заменить картридж")))
			bot.Send(msg)
		}

		if update.Message.Text == "Заменить картридж" {
			customerName := getCustomerNameById(update.Message.Chat.ID)
			text := fmt.Sprintf("Поступила заявка от пользователя %s (ID: %d)", customerName, update.Message.Chat.ID)
			msg := tgbotapi.NewMessage(295415523, text)
			bot.Send(msg)

			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ваша заявка принята! Ожидайте! Скоро будет! 100%!\n\nДа не волнуйтесь Вы!")
			bot.Send(msg)
		}

	}
}
