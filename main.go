package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mymmrac/telego"
)

// handler per /api (per ora dummy)
func apiHandler(w http.ResponseWriter, r *http.Request) {
	// esempio: semplice JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok","message":"API alive OK ;)"}`))
}

// handler per /telegram (dummy, solo 200 OK)
func telegramHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func telegramWebHookHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	botToken := "7915593277:AAFXsBGXij-dsXDNLMb6Orb8zyNXw7fHE7c"

	// Note: Please keep in mind that default logger may expose sensitive information,
	// use in development only
	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Set up a webhook on Telegram side
	_ = bot.SetWebhook(ctx, &telego.SetWebhookParams{
		URL:         "https://crecs-bot.it/telegram",
		SecretToken: bot.SecretToken(),
	})

	// Receive information about webhook
	info, _ := bot.GetWebhookInfo(ctx)
	if err != nil {
		http.Error(w, "failed to get webhook info", http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(info)
	if err != nil {
		http.Error(w, "failed to marshal info", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

func main() {
	// flag --port, default 8080
	port := flag.Int("port", 8080, "Porta HTTP su cui avviare il server")
	flag.Parse()

	mux := http.NewServeMux()

	// --- ROUTE API ---
	// /api oppure /api/... (puoi anche usare /api/ se vuoi un sotto-router)
	mux.HandleFunc("/api", apiHandler)
	// se vuoi che /api/... vada uguale:
	// mux.HandleFunc("/api/", apiHandler)

	// --- ROUTA TELEGRAM WEBHOOK ---
	mux.HandleFunc("/telegram", telegramHandler)

	mux.HandleFunc("/telegramWebHook", telegramWebHookHandler)
	// o /telegram/webhook se vuoi qualcosa di più specifico
	// mux.HandleFunc("/telegram/webhook", telegramHandler)

	// --- STATICI WEBSITE su / ---
	// Tutto ciò che non matcha sopra e parte da / va su website
	websiteDir := http.Dir("website")
	websiteFS := http.FileServer(websiteDir)
	mux.Handle("/", websiteFS)

	// --- STATICI WEBAPP su /webapp ---
	// Serve i file dentro ./webapp quando la URL inizia con /webapp
	webappDir := http.Dir("webapp")

	// 1) Redirect /webapp -> /webapp/
	mux.HandleFunc("/webapp", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/webapp/", http.StatusMovedPermanently)
	})

	// 2) Serve tutto sotto /webapp/
	webappFS := http.StripPrefix("/webapp/", http.FileServer(webappDir))
	mux.Handle("/webapp/", webappFS)

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Server starting on %s\n", addr)

	// controllo banale che le cartelle esistano (solo log, non blocco)
	if _, err := os.Stat("website"); os.IsNotExist(err) {
		log.Println("WARNING: cartella 'website' non trovata")
	}
	if _, err := os.Stat("webapp"); os.IsNotExist(err) {
		log.Println("WARNING: cartella 'webapp' non trovata")
	}

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
