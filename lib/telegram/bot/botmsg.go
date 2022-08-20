package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const (
	sendMessage = "sendMessage"
	TOKEN_ENV   = "TGBOT_TOKEN"
	CHAT_ID_ENV = "CHAT_ID"
)

var TOKEN, BASE_URL, CHAT_ID string

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("botmsg: init: %v", err)
	}
	TOKEN, ok := os.LookupEnv(TOKEN_ENV)
	if !ok {
		log.Fatalf("botmsg: init: no such environment variable: %s", TOKEN_ENV)
	}
	CHAT_ID, ok = os.LookupEnv(CHAT_ID_ENV)
	if !ok {
		log.Fatalf("botmsg: init: no such environment variable: %s", CHAT_ID_ENV)
	}
	BASE_URL = fmt.Sprintf("https://api.telegram.org/bot%s/", TOKEN)
}

func makePayload(content map[string]interface{}) (io.Reader, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(content); err != nil {
		return nil, fmt.Errorf("makePayload: %v", err)
	}
	return buf, nil
}

func SendPlain(text string) error {
	content := map[string]interface{}{"chat_id": CHAT_ID, "text": text}
	return Query(sendMessage, content)
}

func SendHTML(html string) error {
	content := map[string]interface{}{
		"chat_id":    CHAT_ID,
		"text":       html,
		"parse_mode": "HTML",
	}
	return Query(sendMessage, content)
}

func Query(cmd string, content map[string]interface{}) error {
	payload, err := makePayload(content)
	if err != nil {
		return fmt.Errorf("Query: %v", err)
	}
	resp, err := http.Post(BASE_URL+cmd, "application/json", payload)
	if err != nil {
		return fmt.Errorf("Query: %v", err)
	}

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	if res["ok"] == true {
		return nil
	} else {
		return fmt.Errorf("Query: [%d] %s", int(res["error_code"].(float64)), res["description"].(string))
	}
}
