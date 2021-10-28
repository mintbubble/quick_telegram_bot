package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/goombaio/namegenerator"
)

type bResponse struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price,string"`
}

type wallet map[string]float64

var db = map[int]wallet{}

func main() {
	// fmt.Println(mascot.BestMascot())
	// http.HandleFunc("/", hello)
	// http.ListenAndServe("localhost:8000", nil)
	bot, err := tgbotapi.NewBotAPI("MyAwesomeBotToken")
	if err != nil {
		log.Panic(err)
	}

	// bot.Debug = true
	// log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	// rub, err := getRate()
	// 	if err != nil {
	// 		log.Printf(err.Error())
	// 		fmt.Println(err.Error())
	// 		rub = 0
	// 	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		command := strings.Split(update.Message.Text, " ")
		user_id := update.Message.From.ID

		// fmt.Println(command)

		switch strings.ToUpper(command[0]) {
		case "ADD":
			if len(command) != 3 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "wrong format for add command"))
				continue
			}

			_, err := getPrice(strings.ToUpper(command[1]))
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "wrong currency for add command"))
				continue
			}

			money, err := strconv.ParseFloat(command[2], 64)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "wrong format for add command")) //err.Error()
			}

			if _, ok := db[user_id]; !ok {
				db[user_id] = make(wallet)
			}
			db[user_id][strings.ToUpper(command[1])] += money

		case "DEL":
			if len(command) != 2 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "wrong format for del command"))
				log.Printf(err.Error())
				continue
			}
			delete(db[user_id], command[1])

		case "SUB":
			if len(command) != 3 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "wrong format for sub command"))
				log.Printf(err.Error())
				continue
			}
			money, err := strconv.ParseFloat(command[2], 64)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "wrong format for sub command"))
				log.Printf(err.Error())
			}

			if _, ok := db[user_id]; !ok {
				db[user_id] = make(wallet)
			}
			db[user_id][command[1]] -= money

		case "SHOW":
			// fmt.Println(db)
			resp := ""
			for key, value := range db[user_id] {
				usdPrice, err := getPrice(key)
				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "wrong currency for "+key))
					continue
				}
				rubPrice, err := getRate(key, "RUB")
				// if err != nil {
				// 	// bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "wrong currency for " + key + " RUB"))
				// 	rubPrice = 0
				// 	continue
				// }
				eurPrice, err := getRate(key, "EUR")
				// if err != nil {
				// 	// bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "wrong currency for " + key + " EUR"))
				// 	eurPrice = 0
				// 	continue
				// }
				// if err != nil {
				// 	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "wrong currency"))
				// 	continue
				// }
				resp += fmt.Sprintf("%s: %.2f, $%.2f, ₽%.2f, €%.2f\n", key, value, value*usdPrice, value*rubPrice, value*eurPrice)
			}
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, resp))

		default:
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "unknown command"))
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		// msg.ReplyToMessageID = update.Message.MessageID
		// bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text))

	}
}

func getPrice(symbol string) (float64, error) {
	adr := fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%sUSDT", symbol)
	resp, err := http.Get(adr)
	if err != nil {
		return 0, err
	}

	var bRes bResponse
	err = json.NewDecoder(resp.Body).Decode(&bRes)
	if err != nil {
		return 0, err
	}

	if bRes.Symbol == "" {
		return 0, errors.New("wrong currency")
	}

	return bRes.Price, nil
}

func getRate(symbol string, currency string) (float64, error) {
	adr := fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s%s", symbol, currency)
	resp, err := http.Get(adr)
	if err != nil {
		return 0, err
	}

	var bRes bResponse
	err = json.NewDecoder(resp.Body).Decode(&bRes)
	if err != nil {
		return 0, err
	}

	if bRes.Symbol == "" {
		return 0, errors.New("wrong currency")
	}

	return bRes.Price, nil
}

func hello(rw http.ResponseWriter, rq *http.Request) {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)

	name := nameGenerator.Generate()

	fmt.Fprintf(rw, "hello %s", name)
	fmt.Fprintf(rw, "hello %s", rq.URL.Path[1:])

}
