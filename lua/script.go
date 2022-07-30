package lua

const (
	// Lock command line
	lock = `if redis.call("GET", KEYS[1]) == ARGV[1] then
			  redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])
              return "OK"; 
			  else
			return redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])
			  end`

	// UnLock command line
	unLock = `if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
			  else
			return 0
			   end`
)
