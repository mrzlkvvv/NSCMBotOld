package main

import (
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/KirillMerz/NSCMTelegramBot/bot"
	"github.com/KirillMerz/NSCMTelegramBot/nscm"
)

func main() {
	telegramBot := bot.New(
		os.Getenv("TELEGRAM_BOT_TOKEN"),
	)

	go nscm.New(telegramBot).Start()

	telegramBot.Start()
}
