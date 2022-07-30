package lua

import "github.com/go-redis/redis/v8"

type Session interface {
	Client() *redis.Client
}
