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

func getCustomerDataById(id int64) *db.Customer {
	for _, customer := range db.Customers {
		if customer.ID == id {
			return &customer
		}
	}
	return nil
}

func getCustomersDataByAddress(text string) []db.Customer {
	log.Println(text)
	fields := strings.Fields(text)
	log.Println(fields)
	matchedCustomers := []db.Customer{}
Customers:
	for _, customer := range db.Customers {
		log.Println("Адрес, который проверяем:", customer.Address)
		for _, field := range fields {
			if !strings.Contains(customer.Address, field) {
				continue Customers
			}
			log.Println("Вот это слово есть в адресе", field)
		}
		log.Println("Совпадение найдено для адреса:", customer.Address)
		matchedCustomers = append(matchedCustomers, customer)
	}
	return matchedCustomers
}

func addCommentIfExist(text, comment string) string {
	if comment != "" {
		text += "\nПримечание: " + comment
	}
	return text
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

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		//log.Printf("%+v\n", update.Message.Chat)
		if update.Message.Text == "/start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Здравствуйте, я помогу заменить Вам картридж")
			//msg.ReplyToMessageID = update.Message.MessageID
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Заменить картридж")))
			bot.Send(msg)
		}

		if update.Message.Text == "Заменить картридж" {
			customer := getCustomerDataById(update.Message.Chat.ID)
			text := fmt.Sprintf("Поступила новая заявка!\n\nЗаказчик: %s\nID: %d\nАдрес: %s\nТелефон: %s\nПринтер: %s\nКартридж: %s", customer.Name, update.Message.Chat.ID, customer.Address, customer.Phone, customer.Printer, customer.Cartrige)
			text = addCommentIfExist(text, customer.Comment)

			msg := tgbotapi.NewMessage(295415523, text)
			bot.Send(msg)
			//msg = tgbotapi.NewMessage(831891756, text)
			//bot.Send(msg)

			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ваша заявка принята! Ожидайте!\n\nМы свяжемся с Вами в ближайшее время!")
			bot.Send(msg)
		}

		if update.Message.Chat.ID == 295415523 {
			customers := getCustomersDataByAddress(update.Message.Text)
			text := fmt.Sprintf("Найдены следующие совпадения:")
			for _, customer := range customers {
				text += fmt.Sprintf("\n\nЗаказчик: %s\nID: %d\nАдрес: %s\nТелефон: %s\nПринтер: %s\nКартридж: %s", customer.Name, customer.ID, customer.Address, customer.Phone, customer.Printer, customer.Cartrige)
				text = addCommentIfExist(text, customer.Comment)
			}
			msg := tgbotapi.NewMessage(295415523, text)
			bot.Send(msg)
		}
	}
}
