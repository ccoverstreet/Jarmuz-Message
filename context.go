package main

import (
	"bytes"
	"encoding/json"

	"github.com/ccoverstreet/jablkodev"
)

type AppContext struct {
	SendMessage func(string) error
}

func CreateAppContext(conf Config) *AppContext {
	return &AppContext{createSendFunc(conf.BotID)}
}

func createSendFunc(botID string) func(string) error {
	return func(message string) error {
		body := struct {
			BotID string `json:"bot_id"`
			Text  string `json:"text"`
		}{botID, message}

		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return err
		}

		_, err = jablkodev.PostSimple("https://api.groupme.com/v3/bots/post",
			"application/json",
			bytes.NewBuffer(jsonBytes))

		return err
	}
}
