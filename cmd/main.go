package main

import (
	"log"

	"github.com/KirillMerz/NSCMTelegramBot/bot"
)

func main() {
    telegramBot, err := bot.New()
    if err != nil {
        log.Fatalln(err)
    }

    telegramBot.Start()
}
