package main

import (
	"fmt"
	"goapi/bilibili/user"
	"log"
)

func main() {
	cookies, err := user.QRLogin()
	if err != nil {
		log.Fatal(err)
	}
	cookie := user.GetCookie(cookies, "SESSDATA")
	if cookie == nil {
		log.Fatal("no such cookie found: SESSDATA")
	}
	fmt.Println(cookie.Value)
}
