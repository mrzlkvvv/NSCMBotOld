package bot

import (
	"log"
	"os"
	"time"

	"gopkg.in/telebot.v3"
)

func New(botToken string) *telebot.Bot {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:     botToken,
		Poller:    newPoller(botToken),
		ParseMode: "Markdown",
	})

	if err != nil {
		log.Fatalln("BOT: creating error:", err)
	}

	bot.Handle("/help", help)
	bot.Handle("/start", start)
	bot.Handle("/check", check)
	bot.Handle("/unregister", unregister)
	bot.Handle(telebot.OnText, otherHandler)

	return bot
}

func newPoller(botToken string) telebot.Poller {
	webhookUrl, webhookUrlExists := os.LookupEnv("WEBHOOK_URL")
	webhookPort, webhookPortExists := os.LookupEnv("WEBHOOK_PORT")

	if webhookUrlExists && webhookPortExists {
		return newWebhook(webhookUrl, webhookPort, botToken)
	}

	return newLongPoller()
}

func newWebhook(webhookURL, webhookPort, botToken string) *telebot.Webhook {
	tls := &telebot.WebhookTLS{
		Key:  "data/certs/privkey.pem",
		Cert: "data/certs/fullchain.pem",
	}

	webhookEndpoint := &telebot.WebhookEndpoint{
		PublicURL: webhookURL,
	}

	webhook := &telebot.Webhook{
		Listen:         ":" + webhookPort,
		Endpoint:       webhookEndpoint,
		TLS:            tls,
		AllowedUpdates: []string{"message"},
	}

	return webhook
}

func newLongPoller() *telebot.LongPoller {
	return &telebot.LongPoller{
		Timeout:        10 * time.Second,
		AllowedUpdates: []string{"message"},
	}
}
