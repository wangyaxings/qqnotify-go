package main

import (
	"context"
	"log"
	"time"

	"github.com/wangyaxings/qqnotify-go/internal/openidcapture"
	"github.com/wangyaxings/qqnotify-go/internal/smokeenv"
	"github.com/wangyaxings/qqnotify-go/qqnotify"
)

func main() {
	usedPath, err := smokeenv.LoadFirst("examples/smoke/.env.local", ".env.local")
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := qqnotify.LoadConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	log.Printf("listening for a fresh QQ direct message using %s", usedPath)
	log.Print("send a new message to the bot from the QQ account you want to bind")

	openID, content, err := openidcapture.CaptureUserOpenID(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("QQ_USER_OPENID=%s", openID)
	log.Printf("last_message=%s", content)
}
