package bot

import (
	"log"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/telebot.v3"

	"github.com/KirillMerz/NSCMTelegramBot/database"
	"github.com/KirillMerz/NSCMTelegramBot/env"
	"github.com/KirillMerz/NSCMTelegramBot/models"
	"github.com/KirillMerz/NSCMTelegramBot/nscm"
)

var (
	REGISTER_DATA_REGEXP = regexp.MustCompile(`^[А-аЯ-я]{2,20}\s[А-аЯ-я]{2,20}\s[А-аЯ-я]{2,20}\s\d{6}$`)

	db = database.New(env.GetFromEnv("MONGODB_URI"))
)

func start(ctx telebot.Context) error {
	defer logCommand("/start", time.Now(), ctx.Sender())

	var greeting string
	hoursNow := time.Now().Hour()

	switch {
	case (hoursNow > 3) && (hoursNow < 12):
		greeting = "Доброе утро! "
	case (hoursNow > 11) && (hoursNow < 20):
		greeting = "Добрый день! "
	default:
		greeting = "Добрый вечер! "
	}

	err := ctx.Send(greeting + MESSAGE_START)
	if err != nil {
		return err
	}

	return ctx.Send(MESSAGE_HELP)
}

func help(ctx telebot.Context) error {
	defer logCommand("/help", time.Now(), ctx.Sender())
	return ctx.Send(MESSAGE_HELP)
}

func register(ctx telebot.Context) error {
	words := strings.Split(ctx.Message().Text, " ")

	user := models.User{
		ID:         ctx.Sender().ID,
		Lastname:   words[0],
		Name:       words[1],
		SecondName: words[2],
		DocNumber:  words[3],
	}

	results, err := nscm.GetResults(user)
	if err != nil {
		return err
	}

	err = db.RegisterUser(user, results)
	if err != nil {
		return ctx.Send(MESSAGE_DATABASE_ERROR)
	}

	err = ctx.Send(MESSAGE_REGISTER_SUCCESS)
	if err != nil {
		return err
	}

	return sendResults(ctx, results)
}

func unregister(ctx telebot.Context) error {
	defer logCommand("/unregister", time.Now(), ctx.Sender())

	err := db.UnregisterUser(ctx.Sender().ID)
	if err != nil {
		return ctx.Send(MESSAGE_DATABASE_ERROR)
	}

	return ctx.Send(MESSAGE_UNREGISTER_SUCCESS)
}

func check(ctx telebot.Context) error {
	defer logCommand("/check", time.Now(), ctx.Sender())

	results, err := db.GetResults(ctx.Sender().ID)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ctx.Send(MESSAGE_NOT_REGISTERED_ERROR)
		}

		return ctx.Send(MESSAGE_DATABASE_ERROR)
	}

	return sendResults(ctx, results)
}

func otherHandler(ctx telebot.Context) error {
	msg := ctx.Message().Text

	defer logCommand(
		"\""+msg+"\"",
		time.Now(),
		ctx.Sender(),
	)

	if REGISTER_DATA_REGEXP.MatchString(msg) {
		return register(ctx)
	}

	return ctx.Send(MESSAGE_UNKNOWN_COMMAND_ERROR)
}

func sendResults(ctx telebot.Context, results models.Results) error {
	if len(results.List) == 0 {
		return ctx.Send(MESSAGE_RESULTS_NOT_FOUND_ERROR)
	}

	return ctx.Send(nscm.GetResultsMessage(results))
}

func logCommand(funcName string, start time.Time, sender *telebot.User) {
	log.Printf("BOT: %s {%d} (%v)\n", funcName, sender.ID, time.Since(start))
}
