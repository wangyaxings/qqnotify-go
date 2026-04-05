package main

import (
	"context"
	"log"
	"time"

	"github.com/wangyaxings/qqnotify-go/qqnotify"
)

func main() {
	cfg, err := qqnotify.LoadConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	client := qqnotify.NewClient(cfg, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := client.Send(ctx, qqnotify.Notification{
		Title:   "qqnotify-go demo",
		Body:    "Your first QQ notification was sent successfully.",
		Status:  "success",
		Source:  "example-send",
		TraceID: "demo-001",
	}); err != nil {
		log.Fatal(err)
	}

	log.Print("notification sent")
}
