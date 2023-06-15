package nscm

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gopkg.in/telebot.v3"

	"github.com/KirillMerz/NSCMTelegramBot/database"
	"github.com/KirillMerz/NSCMTelegramBot/models"
)

const UpdateInterval = 5 * time.Minute

const MESSAGE_RESULTS_WAS_UPDATED = "Ваши результаты были обновлены\nПосмотрим? (/check)"

type ResultsUpdater struct {
	db  *database.Database
	bot *telebot.Bot
	wg  *sync.WaitGroup
}

func New(bot *telebot.Bot) *ResultsUpdater {
	return &ResultsUpdater{
		db:  database.New(os.Getenv("MONGODB_URI")),
		bot: bot,
		wg:  &sync.WaitGroup{},
	}
}

func (u *ResultsUpdater) Start() {
	log.Println("UPDATER: started...")

	for {
		users, err := u.db.GetAllUsers()
		if err != nil {
			log.Fatalln(err)
		}

		u.processAllUsers(users)

		time.Sleep(UpdateInterval)
	}
}

func (u *ResultsUpdater) processAllUsers(users []models.User) {
	defer func(numOfUsers int, start time.Time) {
		log.Println(fmt.Sprintf("UPDATER: results updated for %d users (%v)", numOfUsers, time.Since(start)))
	}(len(users), time.Now())

	u.wg.Add(len(users))

	for _, user := range users {
		go func(user models.User) {

			defer u.wg.Done()

			err := u.processUser(user)
			if err != nil {
				log.Println(fmt.Sprintf("UPDATER: error processing user {%d}: %v", user.ID, err))
			}

		}(user)
	}

	u.wg.Wait()
}

func (u *ResultsUpdater) processUser(user models.User) error {
	results, err := GetResults(user) // from NSCM
	if err != nil {
		return err
	}

	if len(results.List) == 0 {
		return nil
	}

	oldResults, err := u.db.GetResults(user.ID)
	if err != nil {
		return err
	}

	if len(results.List) == len(oldResults.List) {
		return nil
	}

	err = u.db.ReplaceResults(results)
	if err != nil {
		return err
	}

	log.Printf("UPDATER: results was updated for user {%d}\n", user.ID)

	_, err = u.bot.Send(&telebot.User{ID: user.ID}, MESSAGE_RESULTS_WAS_UPDATED)

	return err
}
