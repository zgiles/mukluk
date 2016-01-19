package main

import (
	"log"
	"time"
	"github.com/garyburd/redigo/redis"
)


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

func (ac appContext) redisgetNodesByField(field string, input string) ([]Node, error) {
		n := []Node{}
		conn := ac.redispool.Get()
		defer conn.Close()
		values, err := redis.Values(conn.Do("SORT", "mssmhpc:mukluk:index:nodes:" + field + ":" + input,
			"BY", "nosort",
			"GET", "#",
			"GET", "mssmhpc:mukluk:nodes:*->hostname",
			"GET", "mssmhpc:mukluk:nodes:*->ipv4address",
			"GET", "mssmhpc:mukluk:nodes:*->macaddress",
			"GET", "mssmhpc:mukluk:nodes:*->os_name",
			"GET", "mssmhpc:mukluk:nodes:*->os_step",
			"GET", "mssmhpc:mukluk:nodes:*->node_type",
			"GET", "mssmhpc:mukluk:nodes:*->oob_type",
			"GET", "mssmhpc:mukluk:nodes:*->heartbeat"))
		if err != nil {
			log.Println(err)
			return n, err
		}
		scanerr := redis.ScanSlice(values, &n)
		if scanerr != nil {
			log.Println(scanerr)
			return n, scanerr
		}
		return n, nil
}

func (ac appContext) redisgetDiscoveredNodesByField(field string, input string) ([]NodesDiscovered, error) {
		n := []NodesDiscovered{}
		conn := ac.redispool.Get()
		defer conn.Close()
		values, err := redis.Values(conn.Do("SORT", "mssmhpc:mukluk:index:discoverednodes:" + field + ":" + input,
			"BY", "nosort",
			"GET", "#",
			"GET", "mssmhpc:mukluk:discoverednodes:*->ipv4address",
			"GET", "mssmhpc:mukluk:discoverednodes:*->macaddress",
			"GET", "mssmhpc:mukluk:discoverednodes:*->surpressed",
			"GET", "mssmhpc:mukluk:discoverednodes:*->enrolled",
			"GET", "mssmhpc:mukluk:discoverednodes:*->checkincount",
			"GET", "mssmhpc:mukluk:discoverednodes:*->heartbeat"))
		if err != nil {
			log.Println(err)
			return n, err
		}
		scanerr := redis.ScanSlice(values, &n)
		if scanerr != nil {
			log.Println(scanerr)
			return n, scanerr
		}
		return n, nil
}
