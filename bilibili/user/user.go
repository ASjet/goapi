package user

import (
	"fmt"
	"goapi/bilibili/api"
	"strconv"
)

const INFO_BASE_URL = "http://api.live.bilibili.com/live_user/v1/Master/info?uid="
const SPACE_BASE_URL = "https://space.bilibili.com/"

type UserInfo struct {
	UID       int
	RoomID    int
	Name      string
	Avatar    string
	SpaceLink string
}

func GetUserInfo(Uid int) (*UserInfo, error) {
	data, err := api.Get(INFO_BASE_URL + strconv.Itoa(Uid))
	if err != nil {
		return nil, fmt.Errorf("GetUserInfo: %v", err)
	}
	info := data["info"].(map[string]interface{})
	uid := int(info["uid"].(float64))
	return &UserInfo{
		UID:       uid,
		RoomID:    int(data["room_id"].(float64)),
		Name:      info["uname"].(string),
		Avatar:    info["face"].(string),
		SpaceLink: SPACE_BASE_URL + strconv.Itoa(uid),
	}, nil
}
