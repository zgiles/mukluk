package main

import (
//	"time"
//	"github.com/garyburd/redigo/redis"
)

func contains(slice []string, item string) bool {
    set := make(map[string]struct{}, len(slice))
    for _, s := range slice {
        set[s] = struct{}{}
    }

    _, ok := set[item]
    return ok
}

/*
// from https://godoc.org/github.com/garyburd/redigo/redis#Pool
func newRedisPool(server, password string) *redis.Pool {
    return &redis.Pool{
        MaxIdle: 3,
        IdleTimeout: 240 * time.Second,
        Dial: func () (redis.Conn, error) {
            c, err := redis.Dial("tcp", server)
            if err != nil {
                return nil, err
            }
						if password != "" {
	            if _, err := c.Do("AUTH", password); err != nil {
	                c.Close()
	                return nil, err
	            }
						}
            return c, err
        },
        TestOnBorrow: func(c redis.Conn, t time.Time) error {
            _, err := c.Do("PING")
            return err
        },
    }
}
*/
