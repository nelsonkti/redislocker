# Redislocker

----
English | [简体中文](README-CN.md)

Redislocker is a set of distributed locks with high availability and high concurrency based on `redis` + `lua`.

Redislocker currently implements `mutex`, with the following features:
* Distributed: support multiple independent machines to run
* Exclusive, features similar to `sync.Mutex`
* Fairness: Follow FIFO
* High performance: avoid the herd effect, etc., reduce the consumption of a lot of cpu, network, etc.
* Guardian coroutine: prevent the task from releasing the lock before the end of the task, and renew the life of the lock
* Prevent lock timeout: prevent the lock from being released for a long time after downtime

-----
Install
```shell
go get -u github.com/nelsonkti/redislocker@latest
```

To start using Redislocker：
```
    var ctx = context.Background()
    var redisClient *redis.Client
    var session *Session
	
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
	
    session, err = NewSession(redisClient)
    locker := RedisLocker(session, key)
    defer locker.Unlock()
    locker.Lock()
```

Benchmark
```shell
goos: darwin
goarch: amd64
pkg: Redislocker
cpu: Intel(R) Core(TM) i5-1038NG7 CPU @ 2.00GHz
BenchmarkLock
BenchmarkLock-8   	   94381	     15908 ns/op	    5381 B/op	      85 allocs/op
PASS
ok  	Redislocker	2.089s
```