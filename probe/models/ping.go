package models

import (
	"fmt"

	"strconv"
	"time"

	"github.com/paulstuart/ping"
)

type PingRes struct {
	Ip   string
	Ping float64
}

func PingProbe(ips []string, limit, timeout int64) []PingRes {
	chLimit := make(chan bool, limit) //控制并发访问量
	chs := make([]chan PingRes, len(ips))

	limitFunc := func(chLimit chan bool, ch chan PingRes, ip string) {
		pingProbe(ip, timeout, ch)
		<-chLimit
	}
	for i, ip := range ips {
		chs[i] = make(chan PingRes, 1)
		chLimit <- true
		go limitFunc(chLimit, chs[i], ip)
	}
	result := []PingRes{}
	for _, ch := range chs {
		res := <-ch
		result = append(result, res)
	}
	return result
}

func pingProbe(address string, timeout int64, ch chan PingRes) {
	var pingRes PingRes
	rtt := pingCheck(address, timeout)
	pingRes.Ip = address
	pingRes.Ping = rtt
	ch <- pingRes
	return
}

func pingCheck(address string, timeout int64) float64 {
	now := time.Now()

	if err := ping.Pinger(address, int(timeout)); err != nil {
		return -1.0
	}

	end := time.Now()
	d := end.Sub(now)

	rttStr := fmt.Sprintf("%.3f", float64(d.Nanoseconds())/1000000.0)
	rtt, _ := strconv.ParseFloat(rttStr, 64)

	return rtt
}
