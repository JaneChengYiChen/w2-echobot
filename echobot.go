package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	checkEnvFile()
	http.HandleFunc("/callback", echoBot)
	if err := http.ListenAndServe(":"+os.Getenv("line_PORT"), nil); err != nil {
		log.Fatal(err)
	}
}

func checkEnvFile() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func echoBot(w http.ResponseWriter, req *http.Request) {
	bot, err := linebot.New(
		os.Getenv("line_CHANNEL_SECRET"),
		os.Getenv("line_CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	events, err := bot.ParseRequest(req)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if message.Text == "city" {
					resp := linebot.NewTextMessage(
						"Hi, which city are you in ?",
					).WithQuickReplies(
						linebot.NewQuickReplyItems(
							linebot.NewQuickReplyButton("https://i.imgur.com/rORVcdH.png", linebot.NewMessageAction("Taipei", "Taipei")),
							linebot.NewQuickReplyButton("https://i.imgur.com/20bgbYT.png", linebot.NewMessageAction("New York", "New York")),
							linebot.NewQuickReplyButton("", linebot.NewLocationAction("Send location")),
						),
					)
					bot.ReplyMessage(event.ReplyToken, resp).Do()
					return
				}
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
					log.Print(err)
				}
			case *linebot.StickerMessage:
				replyMessage := fmt.Sprintf(
					"sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}
