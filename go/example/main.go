package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/magomedcoder/haratsan-sdk/go"
)

func main() {
	client, err := haratsansdk.NewClient("хост", "токен")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	handler := func(ctx context.Context, update *haratsansdk.Update) error {
		content := strings.TrimSpace(update.Content)

		var reply string
		switch {
		case content == "/start":
			reply = "Привет, я пример бота"
		case content == "/time":
			reply = fmt.Sprintf("Сейчас время %s", time.Now())
		default:
			reply = content
		}

		messageId, err := client.SendMessage(ctx, update.FromUserId, reply)
		if err != nil {
			log.Printf("SendMessage: %v", err)
			return err
		}

		log.Printf("Ответ пользователю %d (message_id=%d, sent_id=%d): %q", update.FromUserId, update.MessageId, messageId, reply)

		return nil
	}

	client.RunPolling(ctx, handler)
}
