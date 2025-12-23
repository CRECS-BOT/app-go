package bot

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/CRECS-BOT/app-go/internal/cache"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func LoggingMiddleware() tgbot.Middleware {
	return func(next tgbot.HandlerFunc) tgbot.HandlerFunc {
		return func(ctx context.Context, b *tgbot.Bot, upd *models.Update) {
			if upd.Message != nil && upd.Message.From != nil {
				log.Printf("update: msg from=%d chat=%d type=%s text=%q",
					upd.Message.From.ID,
					upd.Message.Chat.ID,
					upd.Message.Chat.Type,
					upd.Message.Text,
				)
			}
			next(ctx, b, upd)
		}
	}
}

// Redis global rate limit per utente (valido su tutte le repliche k8s)
func RateLimitMiddleware() tgbot.Middleware {
	return func(next tgbot.HandlerFunc) tgbot.HandlerFunc {
		return func(ctx context.Context, b *tgbot.Bot, upd *models.Update) {
			if upd.Message == nil || upd.Message.From == nil {
				next(ctx, b, upd)
				return
			}
			// esempio: applichiamo solo su private
			if upd.Message.Chat.Type != "private" {
				next(ctx, b, upd)
				return
			}

			userID := upd.Message.From.ID
			key := fmt.Sprintf("rl:%d", userID)

			n, err := cache.Rdb.Incr(ctx, key).Result()
			if err != nil {
				// se Redis è giù, non bloccare tutto
				next(ctx, b, upd)
				return
			}
			if n == 1 {
				_ = cache.Rdb.Expire(ctx, key, 10*time.Second).Err()
			}

			if n > 6 {
				_, _ = b.SendMessage(ctx, &tgbot.SendMessageParams{
					ChatID: upd.Message.Chat.ID,
					Text:   "Rate limit: rallenta un attimo 🙏",
				})
				return
			}

			next(ctx, b, upd)
		}
	}
}
