package nodesdb

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"github.com/zgiles/mukluk"
)

type nodesdb struct {
  redispool *redis.Pool
}

// Theory here:
// There will be a mukluk: top key
// then it will be:
//   index:{type}:{key}:value -> {..muids..} for indicies
//   node:{muid} for nodes

func New(redispool *redis.Pool) *nodesdb {
	return &nodesdb{redispool}
}

func (local nodesdb) KVtoMUID(key string, value string) (string, error) {
	a, ae := local.KVtoMUIDs(key, value)
	if ae != nil {
		return "", ae
	}
	switch len(a) {
		case 0:
			return "", errors.New("Zero found in database")
		case 1:
			return a[0], nil
		default:
			return "", errors.New("Key Value returns more than one MUID")
	}
}

func (local nodesdb) KVtoMUIDs(key string, value string) ([]string, error) {
	return local.redisgetMUIDsByField(key, value)
}

func (local nodesdb) MUID(muid string) (mukluk.Node, error) {
	return local.redisgetNodeByMUID(muid)
}

func (local nodesdb) MUIDs(muids []string) ([]mukluk.Node, error) {
	r := []mukluk.Node{}
	for _,i := range muids {
		n, ne := local.MUID(i)
		if ne != nil { return r, ne }
		r = append(r, n)
	}
	return r, nil
}

func (local nodesdb) Update(muid string, key string, value string) (error) {
	return nil
}

func (local nodesdb) redisgetMUIDsByField(key string, value string) ([]string, error) {
		n := []string{}
		conn := local.redispool.Get()
		defer conn.Close()
		values, err := redis.Values(conn.Do("SORT", "mukluk:index:node:" + key + ":" + value,
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

func (nrdb nodesdb) redisgetNodeByMUID(muid string) (mukluk.Node, error) {
		n := mukluk.Node{}
		conn := nrdb.redispool.Get()
		defer conn.Close()
		values, err := redis.Values(conn.Do("HMGET", "mukluk:node:" + muid,
			"uuid",
			"hostname",
			"ipv4address",
			"macaddress",
			"os_name",
			"os_step",
			"node_type",
			"oob_type",
			"heartbeat"))
		if err != nil {
			return n, err
		}
		values, scanerr := redis.Scan(values, &n.Uuid, &n.Hostname, &n.Ipv4address, &n.Macaddress, &n.Os_name, &n.Os_step, &n.Node_type, &n.Oob_type, &n.Heartbeat)
		// scanerr := redis.ScanStruct(values, &n)
		if scanerr != nil {
			return n, scanerr
		}
		return n, nil
}

/*
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
*/
