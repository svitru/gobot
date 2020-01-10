package main

import (
	"log"
	"os"
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func bot(){
    bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
        if err != nil {
                log.Panic(err)
        }

        bot.Debug = false

        log.Printf("Authorized on account %s", bot.Self.UserName)

        u := tgbotapi.NewUpdate(0)
        u.Timeout = 60

        updates, err := bot.GetUpdatesChan(u)

        for update := range updates {

                if update.Message == nil { // ignore any non-Message Updates
                        continue
                }

                msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
                msg.ReplyToMessageID = update.Message.MessageID
                fmt.Println(update.Message.Text)

                if _, err := bot.Send(msg); err != nil {
                        log.Panic(err)
                }
        }
}

func main() {
  go bot()
  keyword := ""
  fmt.Println("type exit to quit")
  for keyword != "exit" {
    fmt.Scan(&keyword) 
  }
}

