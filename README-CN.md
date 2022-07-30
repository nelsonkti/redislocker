# Redislocker

----
[English](README.md) | 简体中文

Redislocker 是在`redis`+ `lua`基础上实现的一套高可用、高并发的分布式锁.

Redislocker 目前实现了`mutex`, 具备以下特点：
* 分布式：支持多台独立的机器运行
* 排它性，特性跟`sync.Mutex`类似
* 公平性：遵循先入先出
* 性能高：避免羊群效应等，减少大量cpu、网络等消耗
* 守护协程：防止任务未结束释放锁，为其锁续命
* 防止锁超时：防止宕机后，导致长时间未释放锁

-----

Redislocker 使用介绍：
下载
```shell
go get -u github.com/nelsonkti/redislocker@latest
```

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