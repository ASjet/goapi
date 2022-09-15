package live

import (
	"fmt"
	"goapi/bilibili/api"
	"goapi/bilibili/user"
	"strconv"
)

const (
	ROOM_INIT_BASEURL   = "http://api.live.bilibili.com/room/v1/Room/room_init?id="
	STATUS_INFO_BASEURL = "https://api.live.bilibili.com/room/v1/Room/get_status_info_by_uids"
	LIVEROOM_BASEURL    = "https://live.bilibili.com/"
	NO_LIVE             = LiveStatus(0)
	LIVING              = LiveStatus(1)
	REPEATING           = LiveStatus(2)
)

type LiveStatus int

func (s LiveStatus) String() string {
	switch s {
	case NO_LIVE:
		return "No live"
	case LIVING:
		return "Living"
	case REPEATING:
		return "Repeating"
	}
	return "Unknown"
}

type LiveInfo struct {
	user.UserInfo
	Title    string
	Status   LiveStatus
	LiveLink string
	Cover    string
}

func parser(data map[string]interface{}) *LiveInfo {
	uid := int(data["uid"].(float64))
	rid := int(data["room_id"].(float64))
	return &LiveInfo{
		UserInfo: user.UserInfo{
			UID:       uid,
			RoomID:    rid,
			Name:      data["uname"].(string),
			Avatar:    data["face"].(string),
			SpaceLink: user.SPACE_BASE_URL + strconv.Itoa(uid),
		},
		Title:    data["title"].(string),
		Status:   LiveStatus(data["live_status"].(float64)),
		LiveLink: LIVEROOM_BASEURL + strconv.Itoa(rid),
		Cover:    data["cover_from_user"].(string),
	}
}

func GetLiveInfoByUIDs(UIDs ...int) (map[int]*LiveInfo, error) {
	data, err := api.Post(STATUS_INFO_BASEURL, map[string]interface{}{
		"uids": UIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("GetLiveInfoByUIDs: %v", err)
	}
	res := make(map[int]*LiveInfo)
	for _, uid := range UIDs {
		res[uid] = parser(data[strconv.Itoa(uid)].(map[string]interface{}))
	}
	return res, nil
}

func GetLiveStatusByRoomID(RoomID int) (LiveStatus, error) {
	data, err := api.Get(ROOM_INIT_BASEURL + strconv.Itoa(RoomID))
	if err != nil {
		return -1, fmt.Errorf("GetLiveStatusByRoomID: %v", err)
	}
	return LiveStatus(data["live_status"].(float64)), nil
}
