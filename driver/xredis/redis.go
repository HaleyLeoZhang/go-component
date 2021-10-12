package xredis

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

func NewPool(c *Config) (*redis.Pool, error) {
	RedisConn := &redis.Pool{
		MaxActive:   c.Pool.MaxActive,
		MaxIdle:     c.Pool.MaxIdle,
		IdleTimeout: c.Pool.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(c.Proto, c.Addr)
			if err != nil {
				return nil, err
			}
			if c.Auth != "" {
				if _, err := conn.Do("AUTH", c.Auth); err != nil {
					conn.Close()
					return nil, err
				}
			}
			return conn, err
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
	}

	return RedisConn, nil
}
