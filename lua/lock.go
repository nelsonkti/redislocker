package lua

import (
	"context"
	"sync"
)

var once sync.Once

func LoadScript(s Session) {
	once.Do(func() {
		LockScript = s.Client().ScriptLoad(context.TODO(), lock).Val()
		UnLockScript = s.Client().ScriptLoad(context.TODO(), unLock).Val()
	})
}
