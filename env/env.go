package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("ENV: .env was not loaded")
	} else {
		log.Println("ENV: .env was loaded")
	}
}

func GetFromEnv(key string) string {
	value, exists := os.LookupEnv(key)

	if !exists {
		log.Fatalf("ENV: '%s' was not found\n", key)
	}

	return value
}
