package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	TelegramToken        string
	TelegramWebhookSecret string
	BotWorkers           int

	MongoURI string
	MongoDB  string

	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

func MustLoadFromEnv() Config {
	get := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("missing env: %s", k)
		}
		return v
	}
	getOpt := func(k, def string) string {
		v := os.Getenv(k)
		if v == "" {
			return def
		}
		return v
	}
	getIntOpt := func(k string, def int) int {
		v := os.Getenv(k)
		if v == "" {
			return def
		}
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("invalid int env %s=%q: %v", k, v, err)
		}
		return n
	}

	return Config{
		TelegramToken:         get("TELEGRAM_BOT_TOKEN"),
		TelegramWebhookSecret: getOpt("TELEGRAM_WEBHOOK_SECRET", ""),
		BotWorkers:            getIntOpt("BOT_WORKERS", 8),

		MongoURI: get("MONGO_URI"),
		MongoDB:  getOpt("MONGO_DB", "mybot"),

		RedisAddr:     get("REDIS_ADDR"),
		RedisPassword: getOpt("REDIS_PASSWORD", ""),
		RedisDB:       getIntOpt("REDIS_DB", 0),
	}
}
