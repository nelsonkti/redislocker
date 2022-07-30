package main

import (
	"context"
)

type Lua struct {
	s *Session
}

func NewLua(s *Session) *Lua {
	return &Lua{s}
}

func (l *Lua) EvalSha(sha1 string, keys []string, args ...interface{}) (interface{}, error) {
	return l.s.Client().EvalSha(context.Background(), sha1, keys, args...).Result()
}
