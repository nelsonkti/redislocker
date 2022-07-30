package redislocker

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/nelsonkti/redislocker/lua"
	"github.com/rs/xid"
	"strconv"
	"time"
)

const (
	defaultTTL            = time.Millisecond * 500
	internalLockLeaseTime = time.Second * 30
	lockSubscribeTimeout  = time.Second * 60
)

var (
	errLockFailed       = errors.New("failed to apply for lock")
	errSubscribeFailed  = errors.New("subscription lock failed")
	errSubscribeTimeout = errors.New("subscription lock timeout")
)

// A Mutex is a redis lock.
type Mutex struct {
	s          *Session
	key        string
	encryptKey string
	uuid       string
	lua        *Lua
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewMutex returns a NewMutex.
func NewMutex(s *Session, key string) *Mutex {
	ctx, cancel := context.WithCancel(context.Background())
	lua.LoadScript(s)
	return &Mutex{
		s:          s,
		key:        key,
		encryptKey: Md5(key),
		lua:        NewLua(s),
		uuid:       xid.New().String(),
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (m *Mutex) Lock() error {
	resp, err := m.tryAcquire()
	if err == redis.Nil {
		resp, err = m.retry()
	} else if err != nil {
		return err
	} else if resp == nil {
		return errLockFailed
	}

	reply, ok := resp.(string)
	if ok && reply == "OK" {
		// daemon lock
		go m.watchDog()
		return nil
	}
	return errLockFailed
}

func (m *Mutex) Unlock() error {
	resp, err := m.lua.EvalSha(lua.UnLockScript, []string{m.key}, []string{m.uuid})
	if err != nil {
		return err
	}

	m.wake()

	_, ok := resp.(int64)
	if !ok {
		return nil
	}

	m.cancel()
	return nil
}

func (m *Mutex) tryAcquire() (interface{}, error) {
	return m.lua.EvalSha(lua.LockScript, []string{m.key}, []string{m.uuid, strconv.Itoa(int(defaultTTL.Milliseconds()))})
}

func (m *Mutex) watchDog() {
	// daemon lock default 30 seconds
	runTime := time.Now().Unix()
	for {
		select {
		case <-m.ctx.Done():
			return
		default:
			now := time.Now().Unix()
			now -= runTime + internalLockLeaseTime.Milliseconds()/1000
			resp := m.s.Client().Get(m.ctx, m.key).Val()
			if now > 0 || resp != "" {
				return
			}
			_ = m.s.Client().PExpire(m.ctx, m.key, internalLockLeaseTime).Err()
			time.Sleep(internalLockLeaseTime / 2)
		}
	}
}

func (m *Mutex) wake() {
	res := m.s.Client().ZRange(m.ctx, m.encryptKey, 0, 0).Val()
	for _, value := range res {
		m.s.Client().Publish(m.ctx, value, value)
		m.s.Client().ZRem(m.ctx, m.encryptKey, value)
	}
}

func (m *Mutex) retry() (interface{}, error) {
	if m.s.Client().ZCard(m.ctx, m.key).Val() > 0 {
		resp, err := m.waitRetry()
		reply, ok := resp.(string)
		if ok && reply == "OK" {
			return resp, err
		}
	}

	err := m.s.Client().ZAdd(m.ctx, m.encryptKey, &redis.Z{Score: float64(time.Now().UnixNano() / 1000), Member: m.uuid}).Err()

	ctx, cancel := context.WithTimeout(context.Background(), lockSubscribeTimeout)
	defer cancel()

	resp, err := m.subscribe(ctx)
	reply, ok := resp.(string)
	if ok && reply == "OK" {
		return resp, err
	}

	return nil, errLockFailed
}

func (m *Mutex) subscribe(ctx context.Context) (interface{}, error) {
	pubSub := m.s.Client().Subscribe(ctx, m.uuid)
	defer pubSub.Close()
	for {
		select {
		case <-ctx.Done():
			_ = pubSub.Unsubscribe(ctx, m.uuid)
			m.s.Client().ZRem(m.ctx, m.encryptKey, m.uuid)
			return nil, errSubscribeTimeout
		case _, ok := <-pubSub.Channel():
			if !ok {
				m.s.Client().ZRem(m.ctx, m.encryptKey, m.uuid)
				return nil, errSubscribeFailed
			}
			return m.waitRetry()
		}
	}
}

func (m *Mutex) waitRetry() (interface{}, error) {
	resp, err := m.tryAcquire()
	if err != nil && err != redis.Nil {
		return resp, err
	}

	reply, ok := resp.(string)
	if ok && reply == "OK" {
		return resp, err
	}
	return nil, nil
}
