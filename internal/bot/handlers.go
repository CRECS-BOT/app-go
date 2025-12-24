package bot

import (
	"context"
	"fmt"
	"strings"

	"github.com/CRECS-BOT/app-go/internal/db"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (d *Dispatcher) HandleStart(ctx context.Context, b *tgbot.Bot, upd *models.Update) {
	m := upd.Message
	if m == nil || m.From == nil {
		return
	}

	_, _ = db.UpsertUserBasic(m.From.ID, m.From.Username, m.From.FirstName)

	txt := "👋 Benvenuto su CRECS-BOT\n" +
		"🎯 Il modo più semplice per creare abbonamenti e condividere contenuti esclusivi su Telegram.\n\n" +
		"Scegli come vuoi iniziare 👇\n"

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{ Text: "👤 Entra come Fan", CallbackData: "button_1" },
				{ Text: "🧑‍🏫 Entra come Creator", CallbackData: "button_2" },
				{ Text: "ℹ️ Come funziona", CallbackData: "button_2" },
				{ Text: "🌍 Lingua", CallbackData: "button_2" },
				{ Text: "🛟 Supporto", CallbackData: "button_2" },
			},
		},
	}

	_, _ = b.SendMessage(ctx, &tgbot.SendMessageParams{ChatID: m.Chat.ID, Text: txt, ReplyMarkup: kb})
}

func (d *Dispatcher) HandleHelp(ctx context.Context, b *tgbot.Bot, upd *models.Update) {
	m := upd.Message
	if m == nil {
		return
	}
	txt := "Help:\n" +
		"- /register avvia un flusso\n" +
		"- /cancel annulla il flusso\n" +
		"- in gruppo rispondo solo se mi menzioni"
	_, _ = b.SendMessage(ctx, &tgbot.SendMessageParams{ChatID: m.Chat.ID, Text: txt})
}

func (d *Dispatcher) HandlePrivateText(ctx context.Context, b *tgbot.Bot, upd *models.Update) {
	m := upd.Message
	if m == nil {
		return
	}

	// Se è un comando non gestito (es. /foo) -> risposta generica
	if strings.HasPrefix(m.Text, "/") {
		_, _ = b.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: m.Chat.ID,
			Text:   "Comando non riconosciuto. /help",
		})
		return
	}

	_, _ = b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: m.Chat.ID,
		Text:   fmt.Sprintf("Hai scritto: %s", m.Text),
	})
}

func (d *Dispatcher) HandleGroupText(ctx context.Context, b *tgbot.Bot, upd *models.Update) {
	m := upd.Message
	if m == nil {
		return
	}

	// Pattern classico: in gruppo rispondi solo se menzionato
	// (nota: dipende da privacy mode del bot)
	botUsername := "" // potresti metterlo in config
	if botUsername != "" && !strings.Contains(m.Text, "@"+botUsername) {
		return
	}

	_, _ = b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: m.Chat.ID,
		Text:   "Sono vivo nel gruppo ✅",
	})
}

func (d *Dispatcher) HandleCallbackAny(ctx context.Context, b *tgbot.Bot, upd *models.Update) {
	cb := upd.CallbackQuery
	if cb == nil {
		return
	}

	// Esempio: data = "flow:confirm" oppure "admin:ban:123"
	data := cb.Data

	// sempre rispondere alla callback per togliere loading
	defer func() {
		_, _ = b.AnswerCallbackQuery(ctx, &tgbot.AnswerCallbackQueryParams{
			CallbackQueryID: cb.ID,
		})
	}()

	_ = data // gestisci routing qui
}
