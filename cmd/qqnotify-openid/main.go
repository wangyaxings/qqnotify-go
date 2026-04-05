package main

import (
	"context"
	"log"
	"net/http"
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

	cfg, err := qqnotify.LoadCaptureConfigFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	accessToken, err := qqnotify.FetchAccessToken(ctx, &http.Client{Timeout: 10 * time.Second}, cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("listening for a fresh QQ direct message using %s", usedPath)
	log.Print("send a new message to the bot from the QQ account you want to bind")

	msg, err := openidcapture.CaptureSingleMessage(ctx, accessToken, func(openidcapture.IncomingC2CMessage) bool {
		return true
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("QQ_USER_OPENID=%s", msg.UserOpenID)
	log.Printf("last_message=%s", msg.Content)
}
