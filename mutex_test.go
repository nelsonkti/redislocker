package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
	"testing"
)

var ctx = context.Background()
var redisClient *redis.Client
var session *Session

func init() {
	redisClient = redis.NewClient(
		&redis.Options{
			Addr:     "0.0.0.0:6379",
			Password: "",
			DB:       0,
		},
	)
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
}

func BenchmarkLock(b *testing.B) {
	var err error
	session, err = NewSession(redisClient)
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		key := xid.New().String()
		go lock(key)
		go lock(key)
		go lock(key)
	}
}

func lock(key string) {
	locker := RedisLocker(session, key)
	defer locker.Unlock()
	locker.Lock()
	// TODO: business logic
}
