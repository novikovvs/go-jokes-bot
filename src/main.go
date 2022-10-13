package main

import (
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/go-co-op/gocron"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/mmcdole/gofeed"
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

	ChannelId, err := strconv.ParseInt(os.Getenv("CHANNEL_ID"), 10, 64)

	if err != nil {
		log.Panic(err)
	}

	bot, err := tgbotapi.NewBotAPI("5487673962:AAFgI0HPIgjF-eugmedPvEGu4NhsbOOgELc")

	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	initCron(ChannelId, bot)

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "help":
			msg.Text = HelpMessage
		case "parse_rss":
			sendDailyJocks(ChannelId, bot)
			continue
		default:
			msg.Text = "Сорри, не вдуплил что тебе надо("
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
func initCron(ChannelId int64, bot *tgbotapi.BotAPI) {
	s := gocron.NewScheduler(time.UTC)
	s.Every(1).Day().At("12:28").Do(func() {
		sendDailyJocks(ChannelId, bot)
	})
	s.StartAsync()
}

func sendDailyJocks(ChannelId int64, bot *tgbotapi.BotAPI) {
	log.Println("Start parsing from RSS")
	msg := tgbotapi.NewMessage(-ChannelId, "")
	msg.Text = "-- Прислано ботом --"
	bot.Send(msg)
	for _, jok := range getJokes() {
		msg.Text = jok
		bot.Send(msg)
	}
	msg.Text = "-- **** --"
	bot.Send(msg)
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
