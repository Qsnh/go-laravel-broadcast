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
	req, err := http.NewRequest("GET", authUrl+"?channel_name="+channel, os.Stdout)
	if err != nil {
		log.Error("init request error.", err)
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.WithField("cookie", cookies).Error("send request error.", err)
		return false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("read auth response content error.", err)
		return false
	}
	responseBody := string(body)
	log.WithField("body", responseBody).Info("auth response content.")
	if responseBody == "true" {
		return true
	}
	return false
}
