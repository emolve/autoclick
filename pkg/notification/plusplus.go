package notification

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

const (
	pushPlusUrl = "http://www.pushplus.plus/send"
)

func Plusplus(token, subject string) error {
	s := struct {
		Token    string `json:"token"`
		Title    string `json:"title"`
		Content  string `json:"content"`
		Template string `json:"template"`
	}{
		Token:    token,
		Title:    subject,
		Content:  time.Now().String() + "打卡成功",
		Template: "html",
	}
	jsonStr, err := json.Marshal(s)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", pushPlusUrl, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return err
}
