package main

import (
	"github.com/KirillMerz/NSCMTelegramBot/bot"
	"github.com/KirillMerz/NSCMTelegramBot/env"
	"github.com/KirillMerz/NSCMTelegramBot/nscm"
)

func main() {
	telegramBot := bot.New(
		env.GetFromEnv("TELEGRAM_BOT_TOKEN"),
	)

	go nscm.New(telegramBot).Start()

	telegramBot.Start()
}
