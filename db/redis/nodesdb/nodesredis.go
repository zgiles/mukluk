package nodesdb

import (
	"log"
	"github.com/garyburd/redigo/redis"
	"github.com/zgiles/mukluk"
)

type nodesdb struct {
  redispool *redis.Pool
}

func New(redispool *redis.Pool) *nodesdb {
	return &nodesdb{redispool}
}

func (nrdb nodesdb) DbSingleKV(field string, input string) (mukluk.Node, error) {
	answer, err := nrdb.redisgetNodesByField(field, input)
	if err != nil {
		return mukluk.Node{}, err
	}
	return answer[0], nil
}

func (nrdb nodesdb) DbMultiKV(field string, input string) ([]mukluk.Node, error) {
	return nrdb.redisgetNodesByField(field, input)
}

func (local nodesdb) MUID(muid string) (mukluk.Node, error) {
	n := mukluk.Node{}
	return n, nil
}

func (nrdb nodesdb) DbUpdateSingleKV(uuid string, key string, value string) (error) {
	return nil
}

func (nrdb nodesdb) redisgetNodesByField(field string, input string) ([]mukluk.Node, error) {
		n := []mukluk.Node{}
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
