package cache

import (
	"errors"
	"fmt"
	"gopkg.in/redis.v5"
	"time"
)

//Memcache struct contains *memcache.Client
type RedisCache struct {
	conn *redis.Client
}

//NewMemcache create new memcache
func NewRedisCache(url string,password string,db int) (Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: url,
		Password: password, // no password set
		DB: db,  // use default DB
	})

	pong, err := client.Ping().Result()
	if err != nil {
		fmt.Println(pong, err)
		return nil, err
	}
	return &RedisCache{client}, nil
}

func NewCache(conn *redis.Client) Cache{
	return &RedisCache{conn}
}

//Get return cached value
func (redis *RedisCache) Get(key string) interface{} {
	if item, err := redis.conn.Get(key).Result(); err == nil {
		return string(item)
	}
	return nil
}

// IsExist check value exists in memcache.
func (redis *RedisCache) IsExist(key string) bool {
	_, err := redis.conn.Get(key).Result()
	if err != nil {
		return false
	}
	return true
}

//Set cached value with key and expire time.
func (redis *RedisCache) Set(key string, val interface{}, timeout time.Duration) error {
	v, ok := val.(string)
	if !ok {
		return errors.New("val must string")
	}

	err := redis.conn.Set(key, v, timeout).Err()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

//Delete delete value in memcache.
func (redis *RedisCache) Delete(key string) error {
	return redis.conn.Del(key).Err()
}
