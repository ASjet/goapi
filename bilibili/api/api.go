package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func Get(URL string, cookies ...*http.Cookie) (map[string]interface{}, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, fmt.Errorf("APIClient.Get: %v", err)
	}
	for _, cookie := range cookies {
		if cookie == nil {
			continue
		}
		req.AddCookie(cookie)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("APIClient.Get: %v", err)
	}
	return Parse(resp)
}

func Post(URL string, Payload map[string]interface{}, cookies ...*http.Cookie) (map[string]interface{}, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(Payload); err != nil {
		return nil, fmt.Errorf("Post: %v", err)
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", URL, buf)
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	req.ParseForm()
	if err != nil {
		return nil, fmt.Errorf("Post: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Post: %v", err)
	}
	return Parse(resp)
}

func Parse(resp *http.Response) (map[string]interface{}, error) {
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
