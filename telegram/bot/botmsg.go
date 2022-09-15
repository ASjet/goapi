package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
