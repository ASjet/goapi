package main

import (
	"bytes"
	"fmt"
	"goapi/bilibili/live"
	"goapi/bilibili/user"
	"goapi/telegram/bot"
	"html/template"
	"log"
	"os"
	"strconv"
	"time"
)

const MessageTemplate = `<u>{{.Name}}</u> <b>Start Living!</b>
<a href="{{.Cover}}">{{.Title}}</a>
Link: {{.LiveLink}}
`

var msg = template.Must(template.New("msg").Parse(MessageTemplate))

func main() {
	// Parse commandline arguments
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <uid> <duration(second)>\n", os.Args[0])
		os.Exit(0)
	}
	uid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	dura, err := time.ParseDuration(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	monitor(uid, dura)
}

func monitor(uid int, interval time.Duration) {
	// Get user info from uid
	ui := getUserInfo(uid)
	log.Printf("Monitoring on %q(uid=%d)'s live room(id=%d)", ui.Name, ui.UID, ui.RoomID)
	log.Printf("Start loop, interval: %s", interval)

	lastStat := live.NO_LIVE
	errCnt := 0

	// Start main loop
	for {
		// Query live status
		stat, err := live.GetLiveStatusByRoomID(ui.RoomID)
		if err != nil || stat == lastStat {
			if err != nil {
				errCnt++
				if errCnt >= 5 {
					log.Printf("err: %v", err)
					errCnt = 0
				}
			}
		} else {
			log.Printf("Live room status updated: %s => %s", lastStat, stat)
			if stat == live.LIVING {
				go func() {
					li := getLiveInfo(uid)
					ui = &li.UserInfo
					log.Printf("%q start living!", li.Name)
					log.Printf("Title: %q", li.Title)
					log.Printf("Address: %s", li.LiveLink)
					sendMessage(li)
				}()
			} else {
				if lastStat == live.LIVING {
					log.Printf("%q stop living!", ui.Name)
				}
			}
			lastStat = stat
		}
		time.Sleep(time.Second * time.Duration(interval))
	}
}

func sendMessage(li *live.LiveInfo) {
	buf := new(bytes.Buffer)
	msg.Execute(buf, li)
	if err := bot.SendHTML(buf.String()); err != nil {
		log.Printf("main: sendMessage: %v", err)
	}
}

func getLiveInfo(uid int) *live.LiveInfo {
	res, err := live.GetLiveInfoByUIDs(uid)
	if err != nil {
		log.Fatalf("main: getLiveInfo: %v", err)
	}
	return res[uid]
}

func getUserInfo(uid int) *user.UserInfo {
	for i := 0; i < 5; i++ {
		res, err := user.GetUserInfo(uid)
		if err != nil {
			log.Printf("main: getUserInfo: %v, retrying...", err)
			time.Sleep(time.Second)
			continue
		}
		return res
	}
	log.Fatalf("can not get user info, uid: %d\n", uid)
	return nil
}
