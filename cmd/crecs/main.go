package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CRECS-BOT/app-go/internal/bot"
	"github.com/CRECS-BOT/app-go/internal/cache"
	"github.com/CRECS-BOT/app-go/internal/config"
	"github.com/CRECS-BOT/app-go/internal/db"
)

func apiHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok","message":"API alive OK ;)"}`))
}

func main() {
	port := flag.Int("port", 8080, "Porta HTTP su cui avviare il server")
	flag.Parse()

	cfg := config.MustLoadFromEnv()
	cfg.PrintConfig()


	// Init Mongo + Redis (shared across replicas -> OK for k8s)
	db.MustInitMongo(cfg.MongoURI, cfg.MongoDB)
	cache.MustInitRedis(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)

	// Init Telegram bot service
	botSvc := bot.MustNewService(cfg)

	// HTTP router
	mux := http.NewServeMux()

	// API
	mux.HandleFunc("/api", apiHandler)

	// Telegram webhook endpoint (IMPORTANT: this must match webhook url)
	mux.Handle("/telegram", botSvc.WebhookHandler())

	// static website on /
	mux.Handle("/", http.FileServer(http.Dir("website")))

	// static webapp on /webapp/
	mux.Handle("/webapp/", http.StripPrefix("/webapp", http.FileServer(http.Dir("webapp"))))

	addr := fmt.Sprintf(":%d", *port)
	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Start bot internal loop for webhook processing
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go botSvc.StartWebhookLoop(ctx)

	log.Printf("HTTP server starting on %s", addr)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutdown...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	_ = server.Shutdown(shutdownCtx)
	botSvc.Close(ctx)

	log.Println("bye")
}

