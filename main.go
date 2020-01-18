package main

import (
	"log"
	"os"
	"fmt"
	"context"
	"strings"

        "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func ConnectDB(client *mongo.Client){
    // Create connect
    err := client.Connect(context.TODO())
    if err != nil {
      log.Fatal(err)
    }

    // Check the connection
    err = client.Ping(context.TODO(), nil)
    if err != nil {
      log.Fatal(err)
    }

    fmt.Println("Connected to MongoDB!")
}

func UpdateStatistic(msg *tgbotapi.Message, client *mongo.Client){
    collectionUsers := client.Database("test").Collection("users")
    collectionChats := client.Database("test").Collection("chats")
    collectionStatistic := client.Database("test").Collection("statistic")

    user := msg.From
    opts := options.Update().SetUpsert(true)        //insert if doc not exist

    filter := bson.D{{"id", user.ID}}
    newitem := bson.D{{"$set", user}}
    _, err := collectionUsers.UpdateOne(context.TODO(), filter, newitem, opts)
    if err != nil {
      log.Fatal(err)
    }

    filter = bson.D{{"id", msg.Chat.ID}}
    newitem = bson.D{{"$set", msg.Chat}}
    _, err = collectionChats.UpdateOne(context.TODO(), filter, newitem, opts)
    if err != nil {
      log.Fatal(err)
    }

    filter = bson.D{{"chat_id", msg.Chat.ID}, {"user_id", msg.From.ID}}
    newitem = bson.D{
                      {"$set", bson.D{{"chat_id", msg.Chat.ID}, {"user_id", msg.From.ID}}},
                      {"$inc", bson.D{{"count", 1}}},
                    }
    _, err = collectionStatistic.UpdateOne(context.TODO(), filter, newitem, opts)
    if err != nil {
      log.Fatal(err)
    }

}

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

    client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
    if err != nil {
      log.Fatal(err)
    }

    ConnectDB(client)

    for update := range updates {

      if update.Message == nil { // ignore any non-Message Updates
        continue
      }

      if strings.Contains(update.Message.Text, "@KangBongSungBot") {
        if update.Message.From.ID == 533587790 {
          msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Моя Настенька! 🤗")
          msg.ReplyToMessageID = update.Message.MessageID

          bot.Send(msg)
        }
      }

    UpdateStatistic(update.Message, client)

    fmt.Printf("%d --- %s --- %d: %s\n", update.Message.From.ID, update.Message.From.FirstName, update.Message.Chat.ID, update.Message.Text)

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

