package bot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/CRECS-BOT/app-go/internal/cache"
	"github.com/CRECS-BOT/app-go/internal/db"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type FSM struct{}

func NewFSM() *FSM { return &FSM{} }

// Avvio flusso: /register
func (d *Dispatcher) HandleRegisterStart(ctx context.Context, b *tgbot.Bot, upd *models.Update) {
	m := upd.Message
	if m == nil || m.From == nil {
		return
	}

	_ = db.SetUserState(m.From.ID, "REG_WAIT_NAME")

	_, _ = b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: m.Chat.ID,
		Text:   "Ok! Dimmi il tuo *nome* (scrivilo qui).",
		ParseMode: "Markdown",
	})
}

// /cancel
func (d *Dispatcher) HandleCancel(ctx context.Context, b *tgbot.Bot, upd *models.Update) {
	m := upd.Message
	if m == nil || m.From == nil {
		return
	}
	_ = db.SetUserState(m.From.ID, "")

	// pulizia chiavi temporanee
	_ = cache.Rdb.Del(ctx, tempKeyName(m.From.ID)).Err()

	_, _ = b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: m.Chat.ID,
		Text:   "Flusso annullato ✅",
	})
}

func (f *FSM) HandleState(ctx context.Context, b *tgbot.Bot, upd *models.Update, u *db.User) {
	m := upd.Message
	if m == nil || m.From == nil {
		return
	}

	text := strings.TrimSpace(m.Text)
	if text == "" {
		return
	}

	switch u.State {
	case "REG_WAIT_NAME":
		// salva nome temporaneo in redis
		_ = cache.Rdb.Set(ctx, tempKeyName(m.From.ID), text, 15*time.Minute).Err()
		_ = db.SetUserState(m.From.ID, "REG_WAIT_EMAIL")

		_, _ = b.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: m.Chat.ID,
			Text:   "Perfetto. Ora dimmi la tua *email*.",
			ParseMode: "Markdown",
		})

	case "REG_WAIT_EMAIL":
		name, _ := cache.Rdb.Get(ctx, tempKeyName(m.From.ID)).Result()
		email := text

		// valida email in modo basico (poi la fai più seria)
		if !strings.Contains(email, "@") {
			_, _ = b.SendMessage(ctx, &tgbot.SendMessageParams{
				ChatID: m.Chat.ID,
				Text:   "Email non valida, riprova.",
			})
			return
		}

		// qui normalmente: salva in Mongo (aggiungi campi al model User)
		// Per esempio: u.FullName = name, u.Email = email, update DB
		_ = db.SetUserState(m.From.ID, "")
		_ = cache.Rdb.Del(ctx, tempKeyName(m.From.ID)).Err()

		_, _ = b.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: m.Chat.ID,
			Text:   fmt.Sprintf("Registrazione completata ✅\nNome: %s\nEmail: %s", name, email),
		})

	default:
		// stato sconosciuto -> reset
		_ = db.SetUserState(m.From.ID, "")
	}
}

func tempKeyName(userID int64) string {
	return fmt.Sprintf("reg:name:%d", userID)
}