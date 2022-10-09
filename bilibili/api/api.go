package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func parse(resp *http.Response) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&res)
	if res["code"] == nil {
		return nil, fmt.Errorf("parse: %v", res)
	}
	if int(res["code"].(float64)) == 0 {
		return res["data"].(map[string]interface{}), nil
	} else {
		return nil, fmt.Errorf("parse: %v", res["msg"])
	}
}

func Get(URL string) (map[string]interface{}, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, fmt.Errorf("Get: %v", err)
	}
	return parse(resp)
}

func Post(URL string, Payload map[string]interface{}) (map[string]interface{}, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(Payload); err != nil {
		return nil, fmt.Errorf("Post: %v", err)
	}
	resp, err := http.Post(URL, "application/json", buf)
	if err != nil {
		return nil, fmt.Errorf("Post: %v", err)
	}
	return parse(resp)
}
