package main

import (
	"bytes"
	"fmt"
	"goapi/lib/bilibili/live"
	"goapi/lib/bilibili/user"
	"goapi/lib/telegram/bot"
	"html/template"
	"log"
	"os"
	"strconv"
	"time"
)

const MessageTemplate = `<b>"{{.Name}}" Start Living!</b>
<a href="{{.Cover}}">{{.Title}}</a>
{{.Address}}
`

var msg = template.Must(template.New("msg").Parse(MessageTemplate))

func main() {
	// Parse commandline arguments
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s uid duration(second)\n", os.Args[0])
		os.Exit(0)
	}
	uid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	dura, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	// Get user info from uid
	ui := getUserInfo(uid)
	log.Printf("Monitoring on %q(uid=%d)'s live room(id=%d)", ui.Name, ui.UID, ui.RoomID)
	log.Printf("Start loop, interval = %d second(s)", dura)

	sent := false

	// Start main loop
	for {
		// Query live status
		stat, err := live.GetLiveStatusByRoomID(ui.RoomID)
		if err != nil {
			log.Printf("main: %v", err)
			continue
		}
		log.Printf("Live room status: %s", stat)
		if stat == live.LIVING {
			if !sent {
				go func() {
					li := getLiveInfo(uid)
					ui = &li.UserInfo
					log.Printf("%q start living!", li.Name)
					log.Printf("Title: %q", li.Title)
					log.Printf("Address: %s", li.Address)
					sendMessage(li)
				}()
				sent = true
			}
		} else {
			if sent {
				log.Printf("%q stop living!", ui.Name)
			}
			sent = false
		}
		time.Sleep(time.Second * time.Duration(dura))
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
	res, err := user.GetUserInfo(uid)
	if err != nil {
		log.Fatalf("main: getUserInfo: %v", err)
	}
	return res
}
