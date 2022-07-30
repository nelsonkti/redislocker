package redislocker

var _ Locker = (*Lock)(nil)

type Locker interface {
	Lock() error
	Unlock() error
}

type Lock struct {
	m *Mutex
}

func (l Lock) Lock() error {
	return l.m.Lock()
}

func (l Lock) Unlock() error {
	return l.m.Unlock()
}

// RedisLocker returns a RedisLock.
func RedisLocker(s *Session, key string) Locker {
	return &Lock{m: NewMutex(s, key)}
}
