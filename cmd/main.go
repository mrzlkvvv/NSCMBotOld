package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/KirillMerz/NSCMBot/bot"
	"github.com/KirillMerz/NSCMBot/nscm"
)

func main() {
	telegramBot := bot.New(
		os.Getenv("TELEGRAM_BOT_TOKEN"),
	)

	resultsUpdater := nscm.New(telegramBot)
	go resultsUpdater.Start()

	telegramBot.Start()
}
