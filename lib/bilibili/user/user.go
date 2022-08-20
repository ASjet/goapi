package user

import (
	"fmt"
	"goapi/lib/bilibili/api"
	"strconv"
)

const BASE_URL = "http://api.live.bilibili.com/live_user/v1/Master/info?uid="

type UserInfo struct {
	UID    int
	RoomID int
	Name   string
	Avatar string
}

func GetUserInfo(Uid int) (*UserInfo, error) {
	data, err := api.Get(BASE_URL + strconv.Itoa(Uid))
	if err != nil {
		return nil, fmt.Errorf("GetUserInfo: %v", err)
	}
	info := data["info"].(map[string]interface{})
	return &UserInfo{
		UID:    int(info["uid"].(float64)),
		RoomID: int(data["room_id"].(float64)),
		Name:   info["uname"].(string),
		Avatar: info["face"].(string),
	}, nil
}
