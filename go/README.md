# Haratsan SDK Go

Клиент для бот-API Haratsan на Go.

Требуется Go 1.26+.

Установка:

```bash
go get github.com/magomedcoder/haratsan-sdk/go
```

```go
package main

import (
	"context"

	haratsansdk "github.com/magomedcoder/haratsan-sdk/go"
)

func main() {
	client, err := haratsansdk.NewClient("хост", "токен")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	handler := func(ctx context.Context, update *haratsansdk.Update) error {
		_, err := client.SendMessage(ctx, update.FromUserId, "Ответ", nil)
		return err
	}

	client.RunPolling(context.Background(), handler)
}
```

### Inline-клавиатура (ReplyMarkup)

```go
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
	{
		{
			Text: "Отмена",
			CallbackData: "cancel",
		},
	},
})
client.SendMessage(ctx, toUserId, "Голосуйте:", markup)
```

### Обработка нажатий кнопок (CallbackQuery)

Когда пользователь нажимает inline-кнопку, бот получает `CallbackQuery`

```go
callbackHandler := func(ctx context.Context, cb *haratsansdk.CallbackQuery) error {
    reply := "Нажато: " + cb.CallbackData
    _, err := client.SendMessage(ctx, cb.FromUserId, reply, nil)
    return err
}

client.RunPolling(ctx, handler, haratsansdk.WithCallbackHandler(callbackHandler))
```

**Пример использования:** [`example/main.go`](example/main.go)