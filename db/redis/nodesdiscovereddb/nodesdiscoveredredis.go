package nodesdiscoveredredis

import (
	"log"
	"github.com/garyburd/redigo/redis"
	"github.com/zgiles/mukluk"
)

type nodesdiscoveredredisdb struct {
  redispool *redis.Pool
}

func New(redispool *redis.Pool) *nodesdiscoveredredisdb {
	return &nodesdiscoveredredisdb{redispool}
}

func (local nodesdiscoveredredisdb) DbSingleKV(field string, input string) (mukluk.NodesDiscovered, error) {
	answer, err := local.redisgetDiscoveredNodesByField(field, input)
	if err != nil {
		return mukluk.NodesDiscovered{}, err
	}
	return answer[0], nil
}

func (local nodesdiscoveredredisdb) DbMultiKV(field string, input string) ([]mukluk.NodesDiscovered, error) {
	return local.redisgetDiscoveredNodesByField(field, input)
}

func (local nodesdiscoveredredisdb) MUID(muid string) (mukluk.NodesDiscovered, error) {
	return mukluk.NodesDiscovered{}, nil
}

func (local nodesdiscoveredredisdb) DbInsert(nd mukluk.NodesDiscovered) (mukluk.NodesDiscovered, error) {
	return mukluk.NodesDiscovered{}, nil
}

func (local nodesdiscoveredredisdb) DbUpdateSingleKV(uuid string, key string, value string) (error) {
	return nil
}

func (local nodesdiscoveredredisdb) redisgetDiscoveredNodesByField(field string, input string) ([]mukluk.NodesDiscovered, error) {
		n := []mukluk.NodesDiscovered{}
		conn := local.redispool.Get()
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
