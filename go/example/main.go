package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	haratsansdk "github.com/magomedcoder/haratsan-sdk/go"
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

		switch {
		case content == "/start":
			_, err := client.SendMessage(ctx, update.FromUserId, "Привет, я пример бота", nil)
			return err

		case content == "/time":
			reply := fmt.Sprintf("Сейчас время %s", time.Now())
			messageId, err := client.SendMessage(ctx, update.FromUserId, reply, nil)
			if err != nil {
				return err
			}
			log.Printf("Ответ пользователю %d (message_id=%d): %q", update.FromUserId, messageId, reply)

			return nil

		case content == "/vote":
			markup := haratsansdk.BuildReplyMarkup([][]haratsansdk.Button{
				{
					{
						Text: "Да",
						CallbackData: "vote_yes",
					},
					{
						Text: "Нет",
						CallbackData: "vote_no",
					},
				},
			})
			_, err := client.SendMessage(ctx, update.FromUserId, "Голосуйте:", markup)
			if err != nil {
				return err
			}
			log.Printf("Отправлено сообщение с кнопками пользователю %d", update.FromUserId)

			return nil

		case content == "/menu":
			markup := haratsansdk.BuildReplyMarkup([][]haratsansdk.Button{
				{
					{
						Text: "Справка",
						CallbackData: "help",
					},
					{
						Text: "Время",
						CallbackData: "time",
					},
				},
				{
					{
						Text: "Отмена",
						CallbackData: "cancel",
					},
				},
			})
			_, err := client.SendMessage(ctx, update.FromUserId, "Выберите действие:", markup)

			return err

		default:
			messageId, err := client.SendMessage(ctx, update.FromUserId, content, nil)
			if err != nil {
				log.Printf("SendMessage: %v", err)
				return err
			}
			log.Printf("Ответ пользователю %d (message_id=%d): %q", update.FromUserId, messageId, content)

			return nil
		}
	}

	client.RunPolling(ctx, handler)
}
