package main

import (
	"log"
	"os"

	"github.com/er28-0652/go-zoomus/zoomus"
)

func main() {
	webhook := os.Getenv("PLUGIN_WEBHOOK_URL")
	token := os.Getenv("PLUGIN_TOKEN")
	msg := &zoomus.Message{
		Title:   os.Getenv("PLUGIN_MSG_TITLE"),
		Summary: os.Getenv("PLUGIN_MSG_SUMMARY"),
		Body:    os.Getenv("PLUGIN_MSG_BODY"),
	}

	zoom, err := zoomus.NewClient(webhook, token)
	if err != nil {
		log.Fatal(err)
	}

	err = zoom.SendMessage(msg)
	if err != nil {
		log.Fatal(err)
	}
}
