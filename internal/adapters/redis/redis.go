package redis

import (
	"context"
	"log"

	goredis "github.com/redis/go-redis/v9"
)

func NewClient(addr, password string) *goredis.Client {
	rdb := goredis.NewClient(&goredis.Options{
		Addr:     addr,
		Password: password, // 👈 this is the only addition
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("redis ping: %v", err)
	}

	log.Printf("connected to redis at %s", addr)
	return rdb
}
