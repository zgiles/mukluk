package osdb

import (
  "errors"
	"github.com/garyburd/redigo/redis"
	"github.com/zgiles/mukluk"
)

type osdb struct {
  redispool *redis.Pool
}

// Theory here:
// There will be a mukluk: top key
// then it will be:
//   os:{os_name}:{os_step} for oses

func New(redispool *redis.Pool) *osdb {
	return &osdb{redispool}
}

func (local osdb) DbSingleNameStep(os_name string, os_step string) (mukluk.Os, error) {
  o := mukluk.Os{}
  conn := local.redispool.Get()
  defer conn.Close()
  // need a bool to check if key exists
  key := "mukluk:os:" + os_name + ":" + os_step
  exists, err := redis.Bool(conn.Do("EXISTS", key))
  if err != nil { return o, err }
  if exists == false { return o, errors.New("Not Found") }
  values, err := redis.Values(conn.Do("HMGET", key,
    "os_name",
    "os_step",
    "boot_mode",
    "boot_kernel",
    "boot_initrd",
    "boot_options",
    "next_step"))
  if err != nil {
    return o, err
  }
  values, scanerr := redis.Scan(values, &o.Os_name, &o.Os_step, &o.Boot_mode, &o.Boot_kernel, &o.Boot_initrd, &o.Boot_options, &o.Next_step)
  if scanerr != nil {
    return o, scanerr
  }
  return o, nil
}
