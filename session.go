package main

import (
	"github.com/go-redis/redis/v8"
)

type Session struct {
	client *redis.Client
}

func NewSession(client *redis.Client) (*Session, error) {
	if err := client.Ping(client.Context()).Err(); err != nil {
		return nil, err
	}
	return &Session{client: client}, nil
}

func (s *Session) Client() *redis.Client {
	return s.client
}
