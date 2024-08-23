package redis

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

var (
	pool      *redis.Pool
	redisHost = "127.0.0.1:6379"
	redisPass = 119742
)

// newRedisPool: create a new pool
func newRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   30,
		IdleTimeout: 300,
		Dial: func() (redis.Conn, error) {
			//open a new connection
			c, err := redis.Dial("tcp", redisHost)
			if err != nil {
				return nil, err
			}
			//auth redis
			if _, err := c.Do("AUTH", redisPass); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
		//test connection
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			//test connection
			_, err := conn.Do("PING")
			return err
		},
	}
}

func init() {
	pool = newRedisPool()
}

func RedisPool() *redis.Pool {
	return pool
}
