package manager

import "sync"

type IRedis interface {
	Start(waitGroup ...*sync.WaitGroup) error
	Stop(waitGroup ...*sync.WaitGroup) error
	Started() bool

	Action(command string, arguments ...string) error

	Quit() (err error)
	Get(key string) (result []byte, err error)
	Type(key string) (result byte, err error)
	Set(key string, arg1 []byte) error
	Save() error
	AllKeys() (result []string, err error)
	Keys(key string) (result []string, err error)
	Exists(key string) (result bool, err error)
	Rename(key, arg1 string) error
	Info() (result map[string]string, err error)
	Ping() error
	Setnx(key string, arg1 []byte) (result bool, err error)
	Getset(key string, arg1 []byte) (result []byte, err error)
	Mget(key string, arg1 []string) (result [][]byte, err error)
	Incr(key string) (result int64, err error)
	Incrby(key string, arg1 int64) (result int64, err error)
	Decr(key string) (result int64, err error)
	Decrby(key string, arg1 int64) (result int64, err error)
	Del(key string) (result bool, err error)
	Randomkey() (result string, err error)
	Renamenx(key string, arg1 string) (result bool, err error)
	Dbsize() (result int64, err error)
	Expire(key string, arg1 int64) (result bool, err error)
	Ttl(key string) (result int64, err error)
	Rpush(key string, arg1 []byte) error
	Lpush(key string, arg1 []byte) error
	Lset(key string, arg1 int64, arg2 []byte) error
	Lrem(key string, arg1 []byte, arg2 int64) (result int64, err error)
	Llen(key string) (result int64, err error)
	Lrange(key string, arg1 int64, arg2 int64) (result [][]byte, err error)
	Ltrim(key string, arg1 int64, arg2 int64) error
	Lindex(key string, arg1 int64) (result []byte, err error)
	Lpop(key string) (result []byte, err error)
	Blpop(key string, timeout int) (result [][]byte, err error)
	Rpop(key string) (result []byte, err error)
	Brpop(key string, timeout int) (result [][]byte, err error)
	Rpoplpush(key string, arg1 string) (result []byte, err error)
	Brpoplpush(key string, arg1 string, timeout int) (result [][]byte, err error)
	Sadd(key string, arg1 []byte) (result bool, err error)
	Srem(key string, arg1 []byte) (result bool, err error)
	Sismember(key string, arg1 []byte) (result bool, err error)
	Smove(key string, arg1 string, arg2 []byte) (result bool, err error)
	Scard(key string) (result int64, err error)
	Sinter(key string, arg1 []string) (result [][]byte, err error)
	Sinterstore(key string, arg1 []string) error
	Sunion(key string, arg1 []string) (result [][]byte, err error)
	Sunionstore(key string, arg1 []string) error
	Sdiff(key string, arg1 []string) (result [][]byte, err error)
	Sdiffstore(key string, arg1 []string) error
	Smembers(key string) (result [][]byte, err error)
	Srandmember(key string) (result []byte, err error)
	Zadd(key string, arg1 float64, arg2 []byte) (result bool, err error)
	Zrem(key string, arg1 []byte) (result bool, err error)
	Zcard(key string) (result int64, err error)
	Zscore(key string, arg1 []byte) (result float64, err error)
	Zrange(key string, arg1 int64, arg2 int64) (result [][]byte, err error)
	Zrevrange(key string, arg1 int64, arg2 int64) (result [][]byte, err error)
	Zrangebyscore(key string, arg1 float64, arg2 float64) (result [][]byte, err error)
	Hget(key string, hashkey string) (result []byte, err error)
	Hset(key string, hashkey string, arg1 []byte) error
	Hgetall(key string) (result [][]byte, err error)
	Flushdb() error
	Flushall() error
	Move(key string, arg1 int64) (result bool, err error)
	Bgsave() error
	Lastsave() (result int64, err error)
	Publish(channel string, message []byte) (recieverCout int64, err error)
}

// RedisConfig ...
type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database int    `json:"database"`
	Password string `json:"password"`
}

// NewRedisConfig...
func NewRedisConfig(host string, port int, database int, password string) *RedisConfig {
	return &RedisConfig{
		Host:     host,
		Port:     port,
		Database: database,
		Password: password,
	}
}

// AddRedis ...
func (manager *Manager) AddRedis(key string, redis IRedis) error {
	manager.redis[key] = redis
	manager.logger.Infof("redis %s added", key)

	return nil
}

// RemoveRedis ...
func (manager *Manager) RemoveRedis(key string) (IRedis, error) {
	redis := manager.redis[key]

	delete(manager.redis, key)
	manager.logger.Infof("redis %s removed", key)

	return redis, nil
}

// GetRedis ...
func (manager *Manager) GetRedis(key string) interface{} {
	if redis, exists := manager.redis[key]; exists {
		return redis
	}
	manager.logger.Infof("redis %s doesn't exist", key)
	return nil
}
