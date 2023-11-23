package bot

import (
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/telebot.v3"

	"github.com/KirillMerz/NSCMTelegramBot/database"
	"github.com/KirillMerz/NSCMTelegramBot/models"
	"github.com/KirillMerz/NSCMTelegramBot/nscm"
)

var (
	REGISTER_DATA_REGEXP = regexp.MustCompile(`^([А-аЯ-яёЁ]{2,20} ){2,3}\d{6}$`)

	db = database.New(os.Getenv("MONGODB_URI"))
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

	err := sendMessage(ctx, greeting+MESSAGE_START)
	if err != nil {
		return err
	}

	return sendMessage(ctx, MESSAGE_HELP)
}

func help(ctx telebot.Context) error {
	defer logCommand("/help", time.Now(), ctx.Sender())
	return sendMessage(ctx, MESSAGE_HELP)
}

func register(ctx telebot.Context) error {
	words := strings.Split(ctx.Message().Text, " ")

	user := models.User{
		ID:        ctx.Sender().ID,
		Lastname:  words[0],
		Name:      words[1],
		DocNumber: words[len(words)-1],
	}

	// Set SecondName, if it exists
	if len(words) == 4 {
		user.SecondName = words[2]
	}

	err := db.RegisterUser(user)
	if err != nil {
		return sendMessage(ctx, MESSAGE_DATABASE_ERROR)
	}

	err = sendMessage(ctx, MESSAGE_REGISTER_SUCCESS)
	if err != nil {
		return err
	}

	results, err := nscm.GetResults(user)
	if err != nil {
		return sendMessage(ctx, MESSAGE_NSCM_IS_A_TEAPOT_ERROR)
	}

	err = db.ReplaceResults(results)
	if err != nil {
		return sendMessage(ctx, MESSAGE_DATABASE_ERROR)
	}

	return sendResults(ctx, results)
}

func unregister(ctx telebot.Context) error {
	defer logCommand("/unregister", time.Now(), ctx.Sender())

	err := db.UnregisterUser(ctx.Sender().ID)
	if err != nil {
		return sendMessage(ctx, MESSAGE_DATABASE_ERROR)
	}

	return sendMessage(ctx, MESSAGE_UNREGISTER_SUCCESS)
}

func check(ctx telebot.Context) error {
	defer logCommand("/check", time.Now(), ctx.Sender())

	results, err := db.GetResults(ctx.Sender().ID)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return sendMessage(ctx, MESSAGE_NOT_REGISTERED_ERROR)
		}

		return sendMessage(ctx, MESSAGE_DATABASE_ERROR)
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

	return sendMessage(ctx, MESSAGE_UNKNOWN_COMMAND_ERROR)
}

func sendResults(ctx telebot.Context, results models.Results) error {
	if len(results.List) == 0 {
		return sendMessage(ctx, MESSAGE_RESULTS_NOT_FOUND_ERROR)
	}

	return sendMessage(ctx, nscm.GetResultsMessage(results))
}

func sendMessage(ctx telebot.Context, message string) error {
	// send message; if bot was blocked by user, then unregister him

	err := ctx.Send(message)

	if err == telebot.ErrBlockedByUser {
		err := db.UnregisterUser(ctx.Sender().ID)
		return err
	}

	return err
}

func logCommand(funcName string, start time.Time, sender *telebot.User) {
	log.Printf("BOT: %s {%d} (%v)\n", funcName, sender.ID, time.Since(start))
}
