package models

import (
	"encoding/json"
	"testing"
)

func Test_PingProbe(t *testing.T) {
	ips := []string{"114.114.114.114", "114.114.115.115", "39.156.69.79"}
	var timeout int64 = 2
	var limit int64 = 1
	res := PingProbe(ips, limit, timeout)
	bs, _ := json.Marshal(res)
	t.Log(string(bs))
}

func Test_UrlProbe(t *testing.T) {
	urls := []string{"http://bbs.ngacn.cc/path/query?aaa=123&bbb=456", "https://www.16.com", "https://www.baidu.com"}
	headers := map[string]string{
		"user-agent": "Mozilla/5.0 (Linux; Android 6.0.1; Moto G (4)) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Mobile Safari/537.36 Edg/87.0.664.66",
	}

	var timeout int64 = 2
	var limit int64 = 3
	res := UrlProbe(urls, headers, limit, timeout)
	bs, _ := json.Marshal(res)
	t.Log(string(bs))
}
