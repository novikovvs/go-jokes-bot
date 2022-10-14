package main

import (
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/mmcdole/gofeed"
	"highjin/bot/backend"
	"log"
	"os"
	"strconv"
	"time"
)

const HelpMessage = "Привет! Я понимаю такие команды\n/parse_rss"

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	AuthKey := os.Getenv("AUTH_TOKEN")

	ChannelId, err := strconv.ParseInt(os.Getenv("CHANNEL_ID"), 10, 64)
	if err != nil {
		log.Panic(err)
	}

	AdminId, err := strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)
	if err != nil {
		log.Panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	initCron(ChannelId, bot, &AuthKey)

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || update.Message.Chat.ID != AdminId {
			continue
		}

		if !update.Message.IsCommand() {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Command() {
		case "help":
			commands, _ := bot.GetMyCommands()
			var text = ""
			for _, command := range commands {
				text += "/" + command.Command + " - " + command.Description + "\n"
			}
			msg.Text = text

		case "health":
			botUser, err := bot.GetMe()
			if err != nil {
				log.Panic(err)
			}

			msg.Text = strconv.FormatInt(botUser.ID, 10) + ":" + botUser.UserName + "/" + botUser.FirstName

		case "key":
			msg.Text = AuthKey

		case "parse_rss":
			sendDailyJocks(ChannelId, bot)
			continue
		default:
			msg.Text = "Сорри, не вдуплил что тебе надо("
		}

		if msg.Text == "" {
			continue
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}

func initCron(ChannelId int64, bot *tgbotapi.BotAPI, AuthKey *string) {
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Day().At("10:28").Do(func() {
		sendDailyJocks(ChannelId, bot)
	})

	if false {
		s.Every(5).Minutes().Do(func() {
			var err error
			AuthKey, err = backend.GetKey(backend.Client(*AuthKey))
			if err != nil {
				log.Println(err)
			}
			log.Println("Update key:", *AuthKey)
		})
	}

	s.StartAsync()
}

func sendDailyJocks(ChannelId int64, bot *tgbotapi.BotAPI) {
	log.Println("Start parsing from RSS")
	msg := tgbotapi.NewMessage(-ChannelId, "")
	msg.Text = "-- Прислано ботом --"
	_, err := bot.Send(msg)
	if err != nil {
		return
	}

	for _, jok := range getJokes() {
		msg.Text = jok
		_, err := bot.Send(msg)
		if err != nil {
			return
		}
	}

	msg.Text = "-- **** --"
	_, err = bot.Send(msg)
	if err != nil {
		return
	}
}

func getJokes() []string {
	return parseRss()
}

func parseRss() []string {
	var jocks []string

	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://www.anekdot.ru/rss/export_j.xml")

	for _, item := range feed.Items {
		jocks = append(jocks, getTextFromHtml(item.GUID))
	}

	return jocks
}

func getTextFromHtml(url string) string {
	var text string

	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{url},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			text = r.HTMLDoc.Find("div.text").Text()
		},
		LogDisabled: true,
	}).Start()

	return text
}
