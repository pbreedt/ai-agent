package telegram

import (
	"context"
	"fmt"
	"net/rpc"
	"os"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// type TelegramBot struct {
// 	aiAgent  *ai.Agent
// 	bot      *bot.Bot
// 	contacts map[string]int
// }

// func New(aiAgent *ai.Agent) *TelegramBot {
// 	t := &TelegramBot{
// 		aiAgent: aiAgent,
// 	}

// 	return t
// }

func StartListner(ctx context.Context, aiClient *rpc.Client) {

	opts := []bot.Option{
		bot.WithDefaultHandler(func(ctx context.Context, b *bot.Bot, update *models.Update) {
			if update == nil || update.Message == nil {
				fmt.Printf("Do nothing with %+v\n", update)
				return
			}
			fmt.Printf("Msg from %s with ChatID %d\n", update.Message.From.Username, update.Message.Chat.ID)
			var response string
			err := aiClient.Call("Agent.RPCRespondToPrompt", update.Message.Text, &response)
			if err != nil {
				response = err.Error()
			}
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   response,
			})
		}),
	}

	b, err := bot.New(os.Getenv("TELEGRAM_BOT_TOKEN"), opts...)
	if err != nil {
		panic(err.Error())
	}

	b.Start(ctx)
}

// func (t *TelegramBot) SendTextMessage(ctx context.Context, to string, message string) {
// 	t.bot.SendMessage(ctx, &bot.SendMessageParams{
// 		ChatID: t.contacts[to],
// 		Text:   message,
// 	})
// }
