package main

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	authUrl = os.Getenv("AUTH_HOST") + os.Getenv("AUTH_PATH")
)

func Authorization(channel string, cookies []*http.Cookie) (bool) {
	client := &http.Client{}
	url := authUrl + "?channel_name=" + channel
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.WithField("cookie", cookies).WithField("url", url).Error("init request error.", err)
		return false
	}
	req.Header.Add("Content-Type", "application/json")
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.WithField("cookie", cookies).WithField("url", url).Error("send request error.", err)
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("read auth response content error.", err)
		return false
	}
	responseBody := string(body)
	log.WithField("cookie", cookies).WithField("status code", resp.StatusCode).WithField("url", url).WithField("content", responseBody).Info("auth response content.")

	if resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}
