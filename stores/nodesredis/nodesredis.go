package nodesredis

import (
	"log"
	"github.com/garyburd/redigo/redis"
	"github.com/zgiles/mukluk/stores/nodes"
)

type nodesredisdb struct {
  redispool *redis.Pool
}

func NewNodesRedis(redispool *redis.Pool) *nodesredisdb {
	return &nodesredisdb{redispool}
}

func (nrdb nodesredisdb) DbSingleKV(field string, input string) (nodes.Node, error) {
	answer, err := nrdb.redisgetNodesByField(field, input)
	if err != nil {
		return nodes.Node{}, err
	}
	return answer[0], nil
}

func (nrdb nodesredisdb) DbMultiKV(field string, input string) ([]nodes.Node, error) {
	return nrdb.redisgetNodesByField(field, input)
}

func (nrdb nodesredisdb) redisgetNodesByField(field string, input string) ([]nodes.Node, error) {
		n := []nodes.Node{}
		conn := nrdb.redispool.Get()
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
