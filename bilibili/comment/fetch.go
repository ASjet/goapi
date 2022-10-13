package comment

import (
	"fmt"
	"goapi/bilibili/api"
	"net/http"
	"strings"

	"github.com/schollz/progressbar/v3"
)

const (
	TYPE_VIDEO   = 1
	TYPE_TOPIC   = 2
	TYPE_LIVE    = 8
	TYPE_ARTICLE = 12
	TYPE_TREND   = 17

	// Number of replies per page
	PAGE_NUM      = 40
	URL_VIEW      = "https://api.bilibili.com/x/web-interface/view"
	URL_REPLY_CNT = "https://api.bilibili.com/x/v2/reply/count"
	URL_REPLY_GET = "https://api.bilibili.com/x/v2/reply"
)

type Comment struct {
	Uid     uint
	Rpid    uint
	Content string
}

func bvid2oid(bvid string) (string, error) {
	data, err := api.Get(URL_VIEW + "?bvid=" + bvid)
	if err != nil {
		return "", fmt.Errorf("bvid2oid: err")
	}
	return data["aid"].(string), nil
}

func oid2bvid(oid string) (string, error) {
	data, err := api.Get(URL_VIEW + "?aid=" + oid)
	if err != nil {
		return "", fmt.Errorf("oid2bvid: err")
	}
	return data["bvid"].(string), nil
}

func getReplyCnt(t, oid int) int {
	params := fmt.Sprintf("?type=%d&oid=%d", t, oid)
	data, err := api.Get(URL_REPLY_CNT + params)
	if err != nil {
		return -1
	}
	return int(data["count"].(float64))
}

func getOnePage(t, oid, ps, pn int, cookies ...*http.Cookie) ([]Comment, error) {
	params := fmt.Sprintf("?type=%d&oid=%d&ps=%d&pn=%d", t, oid, ps, pn)
	data, err := api.Get(URL_REPLY_GET+params, cookies...)
	if err != nil {
		return nil, fmt.Errorf("getOnePage: %v", err)
	}
	res := []Comment{}
	for _, reply := range data["replies"].([]interface{}) {
		r := reply.(map[string]interface{})
		uid := uint(r["mid"].(float64))
		rpid := uint(r["rpid"].(float64))
		content := r["content"].(map[string]interface{})
		message := content["message"].(string)
		res = append(res, Comment{
			Uid:     uid,
			Rpid:    rpid,
			Content: strings.ReplaceAll(message, "\n", ""),
		})
	}
	return res, nil
}

func GetReplies(t, oid int, showProgress bool, cookie ...*http.Cookie) ([]Comment, error) {
	cnt := getReplyCnt(t, oid)
	if cnt == -1 {
		return nil, fmt.Errorf("getReplies: no such oid: %d", oid)
	}
	page_cnt := (cnt / PAGE_NUM) + 1
	res := []Comment{}
	var bar *progressbar.ProgressBar
	if showProgress {
		bar = progressbar.Default(int64(page_cnt))
	}
	for i := 1; i <= page_cnt; i++ {
		batch, _ := getOnePage(t, oid, PAGE_NUM, i, cookie...)
		res = append(res, batch...)
		if showProgress {
			bar.Add(1)
		}
	}
	return res, nil
}

func fetchReplies(t, oid, mod int, ch chan []Comment, cookie ...*http.Cookie) {
	for pn := 1; ; pn++ {
		if pn%mod != 0 {
			continue
		}
		batch, _ := getOnePage(t, oid, PAGE_NUM, pn, cookie...)
		ch <- batch
		if len(batch) == 0 {
			break
		}
	}
	close(ch)
}

func GetReplies2(t, oid, mod int, cookie ...*http.Cookie) chan []Comment {
	ch := make(chan []Comment, 10)
	go fetchReplies(t, oid, mod, ch, cookie...)
	return ch
}
