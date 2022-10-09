package user

import (
	"encoding/json"
	"fmt"
	"goapi/bilibili/api"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mdp/qrterminal/v3"
)

const (
	GET_LOGIN_URL       = "http://passport.bilibili.com/qrcode/getLoginUrl"
	GET_LOGIN_INFO      = "http://passport.bilibili.com/qrcode/getLoginInfo"
	CONTENT_TYPE        = "application/x-www-form-urlencoded"
	SESSION_COOKIE_NAME = "SESSDATA"
)

func dispQR(text string) {
	qrterminal.Generate(text, qrterminal.L, os.Stdout)
}

func getQRcode() (string, string, error) {
	data, err := api.Get(GET_LOGIN_URL)
	if err != nil {
		return "", "", fmt.Errorf("getQRcode: %v", err)
	}
	url, oku := data["url"].(string)
	key, okk := data["oauthKey"].(string)
	if oku && okk {
		return url, key, nil
	}
	return "", "", fmt.Errorf("getQRcode: type assertion failed: %v", data)
}

func QRLogin() (*http.Cookie, error) {
	url, key, err := getQRcode()
	if err != nil {
		return nil, fmt.Errorf("QRLogin: %v", err)
	}
	dispQR(url)
	payload := "oauthKey=" + key
	log.Print("wait for scanning")
	for {
		resp, err := http.Post(GET_LOGIN_INFO, CONTENT_TYPE, strings.NewReader(payload))
		if err != nil {
			return nil, fmt.Errorf("QRLogin: %v", err)
		}
		body := make(map[string]interface{})
		json.NewDecoder(resp.Body).Decode(&body)
		if body["status"].(bool) {
			cookies := resp.Cookies()
			for _, cookie := range cookies {
				if cookie.Name == SESSION_COOKIE_NAME {
					return cookie, nil
				}
			}
			return nil, fmt.Errorf("QRLogin: no cookie %q", SESSION_COOKIE_NAME)
		}
		time.Sleep(time.Second * 2)
	}
}
