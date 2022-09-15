package bot

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	sendMessage = "sendMessage"
	TOKEN_ENV   = "TGBOT_TOKEN"
	CHAT_ID_ENV = "CHAT_ID"
)

var BASE_URL, CHAT_ID string

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("botmsg: init: %v", err)
	}
	token, ok := os.LookupEnv(TOKEN_ENV)
	if !ok {
		log.Fatalf("botmsg: init: no such environment variable: %s", TOKEN_ENV)
	}
	CHAT_ID, ok = os.LookupEnv(CHAT_ID_ENV)
	if !ok {
		log.Fatalf("botmsg: init: no such environment variable: %s", CHAT_ID_ENV)
	}
	BASE_URL = fmt.Sprintf("https://api.telegram.org/bot%s/", token)
}
