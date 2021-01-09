package models

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/toolkits/pkg/logger"
)

type UrlRes struct {
	Url            *url.URL
	Cert           int
	Latency        float64
	HttpStatusCode int
}

// UrlProbe 发起一个 url 拨测
func UrlProbe(rawurls []string, headers map[string]string, limit, timeout int64) []UrlRes {
	chLimit := make(chan bool, limit) //控制并发访问量
	chs := make([]chan UrlRes, len(rawurls))

	limitFunc := func(chLimit chan bool, ch chan UrlRes, rawurl string) {
		urlProbe(rawurl, headers, timeout, ch)
		<-chLimit
	}
	for i, rawurl := range rawurls {
		if err := urlValid(rawurl); err != nil {
			logger.Errorf("url is not valid, ignore: %s", rawurl)
			continue
		}
		chs[i] = make(chan UrlRes, 1)
		chLimit <- true
		go limitFunc(chLimit, chs[i], rawurl)
	}
	result := []UrlRes{}
	for _, ch := range chs {
		res := <-ch
		result = append(result, res)
	}
	return result
}

func urlProbe(rawurl string, headers map[string]string, timeout int64, ch chan UrlRes) {
	urlRes := urlCheck(rawurl, headers, timeout)
	ch <- urlRes
	return
}

func urlCheck(rawurl string, headers map[string]string, timeout int64) (res UrlRes) {
	var urlStruct *url.URL
	urlStruct, _ = url.Parse(rawurl)

	res.Url = urlStruct
	now := time.Now()

	var err error
	var statusCode int
	statusCode, err = httpGet(rawurl, headers, false, timeout)
	if err != nil && strings.Contains(err.Error(), "certificate") {
		now = time.Now()
		res.Cert = -1
		statusCode, err = httpGet(rawurl, headers, true, timeout)
	}
	end := time.Now()
	d := end.Sub(now)

	if err != nil {
		res.Latency = -1.0
		return
	}
	if urlStruct.Scheme == "https" {
		res.Cert = 1
	}

	rttStr := fmt.Sprintf("%.3f", float64(d.Nanoseconds())/1000000.0)
	rtt, _ := strconv.ParseFloat(rttStr, 64)
	res.Latency = rtt
	res.HttpStatusCode = statusCode
	return
}

func httpGet(rawurl string, headers map[string]string, skipCert bool, timeout int64) (statusCode int, err error) {
	req, err := http.NewRequest("GET", rawurl, nil)
	if err != nil {
		return
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Close = true

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipCert},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * time.Duration(timeout),
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	statusCode = resp.StatusCode
	return
}

func urlValid(rawurl string) error {
	_, err := url.Parse(rawurl)
	return err
}
