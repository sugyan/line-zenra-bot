package main

import (
	"github.com/sugyan/go-zenra"
	"github.com/sugyan/line-bot/line"
	"log"
	"net/http"
	"os"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	os.Setenv("HTTPS_PROXY", os.Getenv("FIXIE_URL"))
	bot := &line.Bot{
		ChannelID:     os.Getenv("LINE_CHANNEL_ID"),
		ChannelSecret: os.Getenv("LINE_CHANNEL_SECRET"),
		LineMID:       os.Getenv("LINE_MID"),
	}
	zenrizer := zenra.NewZenrizer()
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		messages, err := bot.ReceiveRequest(req)
		if err != nil {
			log.Print(err)
			return
		}
		sendText := func(to string, message string) {
			result, err := bot.SendText(to, message)
			if err != nil {
				log.Print(err)
				return
			}
			log.Printf("replied: %s", string(result))
		}
		for _, result := range messages.Result {
			if result.Content.ContentType != line.ContentTypeText {
				log.Print("not text message")
				continue
			}
			log.Printf("text message from %s", result.Content.From)
			zenrized := zenrizer.Zenrize(result.Content.Text)
			if zenrized != result.Content.Text {
				go sendText(result.Content.From, zenrized)
			}
		}
	})

	port := os.Getenv("PORT")
	log.Printf("start server :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
