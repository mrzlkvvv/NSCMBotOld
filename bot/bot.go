package bot

import (
	"gopkg.in/telebot.v3"

	"github.com/KirillMerz/NSCMTelegramBot/env"
)

func New() (*telebot.Bot, error) {
    webhook := &telebot.Webhook{
        Listen: ":" + env.GetFromEnv("WEBHOOK_PORT"),
        Endpoint: &telebot.WebhookEndpoint{ PublicURL: env.GetFromEnv("WEBHOOK_ENDPOINT_URL") },
        AllowedUpdates: []string{"message"},
        TLS: &telebot.WebhookTLS{
            Key: "data/certs/privkey.pem",
            Cert: "data/certs/fullchain.pem",
        },
    }

    botPref := telebot.Settings{
        Token: env.GetFromEnv("TELEGRAM_BOT_TOKEN"),
        Poller: webhook,
    }

    bot, err := telebot.NewBot(botPref)
    if err != nil {
        return nil, err
    }

    bot.Handle("/start", start)

    return bot, nil
}
