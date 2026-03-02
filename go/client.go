package haratsansdk

import (
	"context"
	"time"

	"github.com/magomedcoder/haratsan-sdk/go/gen_pb/bot_apipb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const DefaultPollInterval = 2 * time.Second
const DefaultUpdatesLimit = 50

type Update struct {
	MessageId  int64
	FromUserId int64
	Content    string
	CreatedAt  int64
}

type Client struct {
	api   bot_apipb.BotApiServiceClient
	conn  *grpc.ClientConn
	token string
}

func NewClient(addr, token string) (*Client, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(botTokenInterceptor(token)),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		api:   bot_apipb.NewBotApiServiceClient(conn),
		conn:  conn,
		token: token,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetUpdates(ctx context.Context, offset int64, limit int32) ([]*Update, error) {
	resp, err := c.api.GetUpdates(ctx, &bot_apipb.GetUpdatesRequest{
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}

	out := make([]*Update, 0, len(resp.GetUpdates()))
	for _, u := range resp.GetUpdates() {
		out = append(out, &Update{
			MessageId:  u.GetMessageId(),
			FromUserId: u.GetFromUserId(),
			Content:    u.GetContent(),
			CreatedAt:  u.GetCreatedAt(),
		})
	}

	return out, nil
}

func (c *Client) SendMessage(ctx context.Context, toUserId int64, content string, replyMarkup *bot_apipb.ReplyMarkup) (messageId int64, err error) {
	req := &bot_apipb.SendMessageRequest{
		ToUserId: toUserId,
		Content:  content,
	}
	if replyMarkup != nil {
		req.ReplyMarkup = replyMarkup
	}

	resp, err := c.api.SendMessage(ctx, req)
	if err != nil {
		return 0, err
	}

	return resp.GetMessageId(), nil
}

type Button struct {
	Text         string
	CallbackData string
}

func BuildReplyMarkup(rows [][]Button) *bot_apipb.ReplyMarkup {
	protoRows := make([]*bot_apipb.InlineKeyboardRow, 0, len(rows))
	for _, row := range rows {
		buttons := make([]*bot_apipb.InlineKeyboardButton, 0, len(row))
		for _, b := range row {
			buttons = append(buttons, &bot_apipb.InlineKeyboardButton{
				Text:         b.Text,
				CallbackData: b.CallbackData,
			})
		}
		protoRows = append(protoRows, &bot_apipb.InlineKeyboardRow{
			Buttons: buttons,
		})
	}

	return &bot_apipb.ReplyMarkup{
		InlineKeyboard: protoRows,
	}
}

type UpdateHandler func(ctx context.Context, update *Update) error

type PollOption func(*pollConfig)

type pollConfig struct {
	pollInterval time.Duration
	updatesLimit int32
}

func WithPollInterval(d time.Duration) PollOption {
	return func(c *pollConfig) {
		c.pollInterval = d
	}
}

func WithUpdatesLimit(n int32) PollOption {
	return func(c *pollConfig) {
		c.updatesLimit = n
	}
}

func (c *Client) RunPolling(ctx context.Context, handler UpdateHandler, opts ...PollOption) {
	cfg := &pollConfig{
		pollInterval: DefaultPollInterval,
		updatesLimit: DefaultUpdatesLimit,
	}
	for _, o := range opts {
		o(cfg)
	}

	if cfg.pollInterval <= 0 {
		cfg.pollInterval = DefaultPollInterval
	}

	if cfg.updatesLimit <= 0 {
		cfg.updatesLimit = DefaultUpdatesLimit
	}

	var offset int64
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		updates, err := c.GetUpdates(ctx, offset, cfg.updatesLimit)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(cfg.pollInterval):
			}

			continue
		}

		for _, u := range updates {
			if u.MessageId <= 0 {
				continue
			}

			offset = u.MessageId
			if err := handler(ctx, u); err != nil {
				_ = err
			}
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(cfg.pollInterval):
		}
	}
}

func botTokenInterceptor(token string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.AppendToOutgoingContext(ctx, "bot-token", token)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
