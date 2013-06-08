package main

import (
	"encoding/json"
)

type Message struct {
	Author string
	Text   string
}

func (message *Message) toJSON() (string, error) {
	j, err := json.Marshal(message)
	return string(j), err
}

type ChatStorage struct {
	Messages []*Message
}

func (storage *ChatStorage) AddMessage(message *Message) {
	storage.Messages = append(storage.Messages, message)
}
