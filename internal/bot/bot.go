package bot

import (
	"context"
	"log"
	"net/http"

	"github.com/CRECS-BOT/app-go/internal/config"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Service struct {
	cfg config.Config
	b   *tgbot.Bot
	d   *Dispatcher
}

func MustNewService(cfg config.Config) *Service {
	s := &Service{cfg: cfg}
	s.d = NewDispatcher()

	opts := []tgbot.Option{
		tgbot.WithDefaultHandler(s.d.DefaultHandler),
		tgbot.WithWorkers(cfg.BotWorkers),
		tgbot.WithMiddlewares(
			LoggingMiddleware(),
			//RateLimitMiddleware(), // uses Redis
		),
	}

	// Webhook secret validation (recommended)
	if cfg.TelegramWebhookSecret != "" {
		opts = append(opts, tgbot.WithWebhookSecretToken(cfg.TelegramWebhookSecret))
	}

	b, err := tgbot.New(cfg.TelegramToken, opts...)
	if err != nil {
		log.Fatalf("telegram bot init failed: %v", err)
	}
	s.b = b

	// Register specific handlers (commands / callback)
	s.registerHandlers()

	return s
}

func (s *Service) registerHandlers() {
	// Commands
	s.b.RegisterHandler(tgbot.HandlerTypeMessageText, "/start", tgbot.MatchTypeExact, s.d.HandleStart)
	s.b.RegisterHandler(tgbot.HandlerTypeMessageText, "/help", tgbot.MatchTypeExact, s.d.HandleHelp)
	s.b.RegisterHandler(tgbot.HandlerTypeMessageText, "/register", tgbot.MatchTypeExact, s.d.HandleRegisterStart)
	s.b.RegisterHandler(tgbot.HandlerTypeMessageText, "/cancel", tgbot.MatchTypeExact, s.d.HandleCancel)

	// Callback data (inline buttons)
	s.b.RegisterHandler(tgbot.HandlerTypeCallbackQueryData, "", tgbot.MatchTypePrefix, s.d.HandleCallbackAny)

	// You can add: photos, documents, etc with HandlerType...
	_ = models.Update{}
}

func (s *Service) WebhookHandler() http.Handler {
	return s.b.WebhookHandler()
}

func (s *Service) StartWebhookLoop(ctx context.Context) {
	// Required for webhook mode internal processing loop
	s.b.StartWebhook(ctx)
}

func (s *Service) Close(ctx context.Context) {
	s.b.Close(ctx)
}
