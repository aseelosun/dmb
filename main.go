package main

import (
	"doMassageBot/internal/config"
	db2 "doMassageBot/internal/db"
	"doMassageBot/internal/telegram"
	"fmt"
	_ "github.com/lib/pq"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"log"
	"net/http"
)

func main() {

	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("ok"))
		})

		_ = http.ListenAndServe(":8080", nil)
	}()

	conf, err := config.LoadConfiguration("config.json")
	if err != nil {
		fmt.Printf("error config, %s", err)
		panic(1)
	}
	db, err := db2.ConnectingToDb(&conf)
	if err != nil {
		fmt.Printf("error ConnectingToDb, %s", err)
		panic(1)
	}

	botApi, err := tgbotapi.NewBotAPI(conf.TelegramBotToken)
	if err != nil {
		fmt.Println(err)
		panic(1)
	}
	bot := telegram.NewBot(botApi)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = conf.UpdateTimeout

	if err := bot.Start(db, &conf); err != nil {
		log.Fatal(err)
	}
}
