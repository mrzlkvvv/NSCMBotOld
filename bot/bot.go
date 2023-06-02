package bot

import (
	"log"

	"gopkg.in/telebot.v3"
)

func New(webhookURL, webhookPort, botToken string) *telebot.Bot {
	return newAPI(webhookURL, webhookPort, botToken)
}

func newAPI(webhookURL, webhookPort, botToken string) *telebot.Bot {
	bot, err := telebot.NewBot(telebot.Settings{
		Token:  botToken,
		Poller: newWebhook(webhookURL, webhookPort, botToken),
	})

	if err != nil {
		log.Fatalln(err)
	}

	bot.Handle("/help", help)
	bot.Handle("/start", start)
	bot.Handle("/check", check)
	bot.Handle("/unregister", unregister)
	bot.Handle(telebot.OnText, otherHandler)

	return bot
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
