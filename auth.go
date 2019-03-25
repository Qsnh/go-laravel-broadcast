package main

import (
	"io/ioutil"
	"net/http"
	"os"
)

var (
	authUrl = GetEnv("AUTH_HOST") + GetEnv("AUTH_PATH")
)

func Authorization(channel string, cookies []*http.Cookie) (bool) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", authUrl+"?channel_name="+channel, os.Stdout)
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return false
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	responseBody := string(body)
	if responseBody == "true" {
		return true
	}
	return false
}
