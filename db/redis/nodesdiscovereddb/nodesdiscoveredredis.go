package nodesdiscovereddb

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"github.com/zgiles/mukluk"
)

type nodesdiscovereddb struct {
  redispool *redis.Pool
}

func New(redispool *redis.Pool) *nodesdiscovereddb {
	return &nodesdiscovereddb{redispool}
}

func (local nodesdiscovereddb) KVtoMUID(key string, value string) (string, error) {
	a, ae := local.KVtoMUIDs(key, value)
	if ae != nil {
		return "", ae
	}
	switch len(a) {
		case 1:
			return a[0], nil
		default:
			return "", errors.New("Key Value returns more than one MUID")
	}
}

func (local nodesdiscovereddb) KVtoMUIDs(key string, value string) ([]string, error) {
	return local.redisgetMUIDsByField(key, value)
}

func (local nodesdiscovereddb) MUID(muid string) (mukluk.NodesDiscovered, error) {
	return local.redisgetNodeDiscoveredByMUID(muid)
}

func (local nodesdiscovereddb) MUIDs(muids []string) ([]mukluk.NodesDiscovered, error) {
	r := []mukluk.NodesDiscovered{}
	for _,i := range muids {
		n, ne := local.MUID(i)
		if ne != nil { return r, ne }
		r = append(r, n)
	}
	return r, nil
}

func (local nodesdiscovereddb) Update(muid string, key string, value string) (error) {
	return nil
}

func (local nodesdiscovereddb) Insert(nd mukluk.NodesDiscovered) (mukluk.NodesDiscovered, error) {
	return nd, nil
}


func (local nodesdiscovereddb) redisgetMUIDsByField(field string, input string) ([]string, error) {
		n := []string{}
		conn := local.redispool.Get()
		defer conn.Close()
		values, err := redis.Values(conn.Do("SORT", "mssmhpc:mukluk:index:nodediscovered:" + field + ":" + input,
			"BY", "nosort",
			"GET", "#"))
		if err != nil {
			return n, err
		}
		scanerr := redis.ScanSlice(values, &n)
		if scanerr != nil {
			return n, scanerr
		}
		return n, nil
}

func (local nodesdiscovereddb) redisgetNodeDiscoveredByMUID(muid string) (mukluk.NodesDiscovered, error) {
		n := mukluk.NodesDiscovered{}
		conn := local.redispool.Get()
		defer conn.Close()
		values, err := redis.Values(conn.Do("GET", "mukluk:nodediscovered:" + muid))
		if err != nil {
			return n, err
		}
		scanerr := redis.ScanStruct(values, &n)
		if scanerr != nil {
			return n, scanerr
		}
		return n, nil
}


/*
func (local nodesdiscovereddb) DbSingleKV(field string, input string) (mukluk.NodesDiscovered, error) {
	answer, err := local.redisgetDiscoveredNodesByField(field, input)
	if err != nil {
		return mukluk.NodesDiscovered{}, err
	}
	return answer[0], nil
}

func (local nodesdiscovereddb) DbMultiKV(field string, input string) ([]mukluk.NodesDiscovered, error) {
	return local.redisgetDiscoveredNodesByField(field, input)
}

func (local nodesdiscovereddb) redisgetDiscoveredNodesByField(field string, input string) ([]mukluk.NodesDiscovered, error) {
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

*/
