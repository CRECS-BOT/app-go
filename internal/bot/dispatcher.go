package bot

import (
	"context"
	"strings"

	"github.com/CRECS-BOT/app-go/internal/db"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Dispatcher struct {
	fsm *FSM
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		fsm: NewFSM(),
	}
}

// DefaultHandler = catch-all (messages not matched by RegisterHandler)
func (d *Dispatcher) DefaultHandler(ctx context.Context, b *tgbot.Bot, upd *models.Update) {
	// 1) callback?
	if upd.CallbackQuery != nil {
		d.HandleCallbackAny(ctx, b, upd)
		return
	}

	// 2) message?
	if upd.Message == nil {
		return
	}

	msg := upd.Message
	chatType := msg.Chat.Type

	// Ensure user exists / update basic profile (for private)
	if msg.From != nil {
		_, _ = db.UpsertUserBasic(msg.From.ID, msg.From.Username, msg.From.FirstName)
	}

	// 3) if user has active FSM state, handle it first (private only, di solito)
	if msg.From != nil && chatType == "private" && msg.Text != "" {
		u, _ := db.FindUserByTelegramID(msg.From.ID)
		if u != nil && u.State != "" && !strings.HasPrefix(msg.Text, "/") {
			d.fsm.HandleState(ctx, b, upd, u)
			return
		}
	}

	// 4) dispatch by chat type + update type
	if msg.Text != "" {
		switch chatType {
		case "private":
			d.HandlePrivateText(ctx, b, upd)
		case "group", "supergroup":
			d.HandleGroupText(ctx, b, upd)
		case "channel":
			// channel post (depends on bot permissions)
			// optional: ignore
		}
	}
}
