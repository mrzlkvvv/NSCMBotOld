package main

import (
	"github.com/KirillMerz/NSCMTelegramBot/bot"
	"github.com/KirillMerz/NSCMTelegramBot/env"
)

func main() {
	telegramBot := bot.New(
		env.GetFromEnv("WEBHOOK_URL"),
		env.GetFromEnv("WEBHOOK_PORT"),
		env.GetFromEnv("TELEGRAM_BOT_TOKEN"),
	)

	telegramBot.Start()
}
