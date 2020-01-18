package main

import (
	"log"
	"os"
	"fmt"
	"context"

        "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

    // Create client
    client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
    if err != nil {
      log.Fatal(err)
    }

    // Create connect
    err = client.Connect(context.TODO())
    if err != nil {
      log.Fatal(err)
    }

    // Check the connection
    err = client.Ping(context.TODO(), nil)
    if err != nil {
      log.Fatal(err)
    }

    fmt.Println("Connected to MongoDB!")
    collection := client.Database("test").Collection("chats")

        for update := range updates {

                if update.Message == nil { // ignore any non-Message Updates
                        continue
                }

                msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
                msg.ReplyToMessageID = update.Message.MessageID
                fmt.Println(update.Message.Text)
		user := update.Message.From
		opts := options.Update().SetUpsert(true)
		filter := bson.D{{"id", user.ID}}
		newitem := bson.D{{"$set", user, }}
		updateResult, err := collection.UpdateOne(context.TODO(), filter, newitem, opts)
                if err != nil {
                  log.Fatal(err)
                }

		fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
		fmt.Printf("%d --- %s\n", user.ID, user.FirstName)

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

