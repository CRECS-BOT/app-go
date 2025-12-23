package cache

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

func MustInitRedis(addr, password string, db int) {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := Rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("redis ping failed: %v", err)
	}
	log.Println("redis: connected")
}
