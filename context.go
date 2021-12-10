package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type AppContext struct {
	SendMessage func(string) error
}

func CreateAppContext(conf Config) *AppContext {
	return &AppContext{createSendFunc(conf.BotID)}
}

func createSendFunc(botID string) func(string) error {
	return func(message string) error {
		type bodyFormat struct {
			BotID string `json:"bot_id"`
			Text  string `json:"text"`
		}

		body := bodyFormat{botID, message}

		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return err
		}

		req, err := http.NewRequest("POST", "https://api.groupme.com/v3/bots/post", bytes.NewBuffer(jsonBytes))
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode < 200 || res.StatusCode >= 400 {
			return fmt.Errorf("Status code %d from GroupMe", res.StatusCode)
		}

		return nil
	}
}
