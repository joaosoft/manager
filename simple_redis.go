package manager

import (
	"fmt"
	"github.com/joaosoft/logger"
	"reflect"

	"strings"

	"sync"

	"github.com/alphazero/Go-Redis"
)

// SimpleRedis ...
type SimpleRedis struct {
	client  redis.Client
	config  *RedisConfig
	logger logger.ILogger
	started bool
}

// NewSimpleRedis ...
func (manager *Manager) NewSimpleRedis(config *RedisConfig) IRedis {
	return &SimpleRedis{
		config: config,
		logger:manager.logger,
	}
}

// Start ...
func (redis *SimpleRedis) Start(waitGroup ...*sync.WaitGroup) error {
	var wg *sync.WaitGroup

	if len(waitGroup) == 0 {
		wg = &sync.WaitGroup{}
		wg.Add(1)
	} else {
		wg = waitGroup[0]
	}

	defer wg.Done()

	if redis.started {
		return nil
	}

	if conn, err := redis.config.Connect(); err != nil {
		redis.logger.Error(err)
		return err
	} else {
		redis.client = conn
		redis.started = true
	}
	return nil
}

// Stop ...
func (redis *SimpleRedis) Stop(waitGroup ...*sync.WaitGroup) error {
	var wg *sync.WaitGroup

	if len(waitGroup) == 0 {
		wg = &sync.WaitGroup{}
		wg.Add(1)
	} else {
		wg = waitGroup[0]
	}

	defer wg.Done()

	if !redis.started {
		return nil
	}

	if err := redis.client.Quit(); err != nil {
		return err
	}

	redis.started = false

	return nil
}

// Started ...
func (redis *SimpleRedis) Started() bool {
	return redis.started
}

// Action ...
func (redis *SimpleRedis) Action(command string, arguments ...string) error {
	inputs := make([]reflect.Value, len(arguments))
	for i, arg := range arguments {
		inputs[i] = reflect.ValueOf(arg)
	}

	chars := strings.SplitN(command, "", 1)
	method := fmt.Sprintf("%s%s", strings.ToUpper(chars[0]), strings.ToLower(chars[1]))

	result := reflect.ValueOf(redis.client).MethodByName(method).Call(inputs)
	if result != nil && len(result) > 0 && !result[0].IsNil() {
		return fmt.Errorf(result[0].String())
	}

	return nil
}

func (redis *SimpleRedis) Quit() error {
	return redis.client.Quit()
}

func (redis *SimpleRedis) Get(key string) ([]byte, error) {
	return redis.client.Get(key)

}

func (redis *SimpleRedis) Type(key string) (byte, error) {
	res, err := redis.client.Type(key)
	return byte(res), err
}

func (redis *SimpleRedis) Set(key string, arg1 []byte) error {
	return redis.client.Set(key, arg1)
}

func (redis *SimpleRedis) Save() error {
	return redis.client.Save()
}

func (redis *SimpleRedis) AllKeys() ([]string, error) {
	return redis.client.AllKeys()
}

func (redis *SimpleRedis) Keys(key string) ([]string, error) {
	return redis.client.Keys(key)
}

func (redis *SimpleRedis) Exists(key string) (bool, error) {
	return redis.client.Exists(key)
}

func (redis *SimpleRedis) Rename(key, arg1 string) error {
	return redis.client.Rename(key, arg1)
}

func (redis *SimpleRedis) Info() (map[string]string, error) {
	return redis.client.Info()
}

func (redis *SimpleRedis) Ping() error {
	return redis.client.Ping()
}

func (redis *SimpleRedis) Setnx(key string, arg1 []byte) (bool, error) {
	return redis.client.Setnx(key, arg1)
}

func (redis *SimpleRedis) Getset(key string, arg1 []byte) ([]byte, error) {
	return redis.client.Getset(key, arg1)
}

func (redis *SimpleRedis) Mget(key string, arg1 []string) ([][]byte, error) {
	return redis.client.Mget(key, arg1)
}

func (redis *SimpleRedis) Incr(key string) (int64, error) {
	return redis.client.Incr(key)
}

func (redis *SimpleRedis) Incrby(key string, arg1 int64) (int64, error) {
	return redis.client.Incrby(key, arg1)
}

func (redis *SimpleRedis) Decr(key string) (int64, error) {
	return redis.client.Decr(key)
}

func (redis *SimpleRedis) Decrby(key string, arg1 int64) (int64, error) {
	return redis.client.Decrby(key, arg1)
}

func (redis *SimpleRedis) Del(key string) (bool, error) {
	return redis.client.Del(key)
}

func (redis *SimpleRedis) Randomkey() (string, error) {
	return redis.client.Randomkey()
}

func (redis *SimpleRedis) Renamenx(key string, arg1 string) (bool, error) {
	return redis.client.Renamenx(key, arg1)
}

func (redis *SimpleRedis) Dbsize() (result int64, err error) {
	return redis.client.Dbsize()
}

func (redis *SimpleRedis) Expire(key string, arg1 int64) (bool, error) {
	return redis.client.Expire(key, arg1)
}

func (redis *SimpleRedis) Ttl(key string) (int64, error) {
	return redis.client.Ttl(key)
}

func (redis *SimpleRedis) Rpush(key string, arg1 []byte) error {
	return redis.client.Rpush(key, arg1)
}

func (redis *SimpleRedis) Lpush(key string, arg1 []byte) error {
	return redis.client.Lpush(key, arg1)
}

func (redis *SimpleRedis) Lset(key string, arg1 int64, arg2 []byte) error {
	return redis.client.Lset(key, arg1, arg2)
}

func (redis *SimpleRedis) Lrem(key string, arg1 []byte, arg2 int64) (int64, error) {
	return redis.client.Lrem(key, arg1, arg2)
}

func (redis *SimpleRedis) Llen(key string) (int64, error) {
	return redis.client.Llen(key)
}

func (redis *SimpleRedis) Lrange(key string, arg1 int64, arg2 int64) ([][]byte, error) {
	return redis.client.Lrange(key, arg1, arg2)
}

func (redis *SimpleRedis) Ltrim(key string, arg1 int64, arg2 int64) error {
	return redis.client.Ltrim(key, arg1, arg2)
}

func (redis *SimpleRedis) Lindex(key string, arg1 int64) ([]byte, error) {
	return redis.client.Lindex(key, arg1)
}

func (redis *SimpleRedis) Lpop(key string) ([]byte, error) {
	return redis.client.Lpop(key)
}

func (redis *SimpleRedis) Blpop(key string, timeout int) ([][]byte, error) {
	return redis.client.Blpop(key, timeout)
}

func (redis *SimpleRedis) Rpop(key string) ([]byte, error) {
	return redis.client.Rpop(key)
}

func (redis *SimpleRedis) Brpop(key string, timeout int) ([][]byte, error) {
	return redis.client.Brpop(key, timeout)
}

func (redis *SimpleRedis) Rpoplpush(key string, arg1 string) ([]byte, error) {
	return redis.client.Rpoplpush(key, arg1)
}

func (redis *SimpleRedis) Brpoplpush(key string, arg1 string, timeout int) ([][]byte, error) {
	return redis.client.Brpoplpush(key, arg1, timeout)
}
func (redis *SimpleRedis) Sadd(key string, arg1 []byte) (bool, error) {
	return redis.client.Sadd(key, arg1)
}

func (redis *SimpleRedis) Srem(key string, arg1 []byte) (bool, error) {
	return redis.client.Srem(key, arg1)
}

func (redis *SimpleRedis) Sismember(key string, arg1 []byte) (bool, error) {
	return redis.client.Sismember(key, arg1)
}

func (redis *SimpleRedis) Smove(key string, arg1 string, arg2 []byte) (bool, error) {
	return redis.client.Smove(key, arg1, arg2)
}

func (redis *SimpleRedis) Scard(key string) (int64, error) {
	return redis.client.Scard(key)
}

func (redis *SimpleRedis) Sinter(key string, arg1 []string) ([][]byte, error) {
	return redis.client.Sinter(key, arg1)
}

func (redis *SimpleRedis) Sinterstore(key string, arg1 []string) error {
	return redis.client.Sinterstore(key, arg1)
}

func (redis *SimpleRedis) Sunion(key string, arg1 []string) ([][]byte, error) {
	return redis.client.Sunion(key, arg1)
}

func (redis *SimpleRedis) Sunionstore(key string, arg1 []string) error {
	return redis.client.Sunionstore(key, arg1)
}

func (redis *SimpleRedis) Sdiff(key string, arg1 []string) ([][]byte, error) {
	return redis.client.Sdiff(key, arg1)
}

func (redis *SimpleRedis) Sdiffstore(key string, arg1 []string) error {
	return redis.client.Sdiffstore(key, arg1)
}

func (redis *SimpleRedis) Smembers(key string) ([][]byte, error) {
	return redis.client.Smembers(key)
}

func (redis *SimpleRedis) Srandmember(key string) ([]byte, error) {
	return redis.client.Srandmember(key)
}

func (redis *SimpleRedis) Zadd(key string, arg1 float64, arg2 []byte) (bool, error) {
	return redis.client.Zadd(key, arg1, arg2)
}

func (redis *SimpleRedis) Zrem(key string, arg1 []byte) (bool, error) {
	return redis.client.Zrem(key, arg1)
}

func (redis *SimpleRedis) Zcard(key string) (int64, error) {
	return redis.client.Zcard(key)
}

func (redis *SimpleRedis) Zscore(key string, arg1 []byte) (float64, error) {
	return redis.client.Zscore(key, arg1)
}

func (redis *SimpleRedis) Zrange(key string, arg1 int64, arg2 int64) ([][]byte, error) {
	return redis.client.Zrange(key, arg1, arg2)
}

func (redis *SimpleRedis) Zrevrange(key string, arg1 int64, arg2 int64) ([][]byte, error) {
	return redis.client.Zrevrange(key, arg1, arg2)
}

func (redis *SimpleRedis) Zrangebyscore(key string, arg1 float64, arg2 float64) ([][]byte, error) {
	return redis.client.Zrangebyscore(key, arg1, arg2)
}

func (redis *SimpleRedis) Hget(key string, hashkey string) ([]byte, error) {
	return redis.client.Hget(key, hashkey)
}

func (redis *SimpleRedis) Hset(key string, hashkey string, arg1 []byte) error {
	return redis.client.Hset(key, hashkey, arg1)
}

func (redis *SimpleRedis) Hgetall(key string) ([][]byte, error) {
	return redis.client.Hgetall(key)
}

func (redis *SimpleRedis) Flushdb() error {
	return redis.client.Flushdb()
}

func (redis *SimpleRedis) Flushall() error {
	return redis.client.Flushall()
}

func (redis *SimpleRedis) Move(key string, arg1 int64) (bool, error) {
	return redis.client.Move(key, arg1)
}

func (redis *SimpleRedis) Bgsave() error {
	return redis.client.Bgsave()
}

func (redis *SimpleRedis) Lastsave() (int64, error) {
	return redis.client.Lastsave()
}

func (redis *SimpleRedis) Publish(channel string, message []byte) (int64, error) {
	return redis.client.Publish(channel, message)
}
