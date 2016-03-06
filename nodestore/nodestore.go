package nodestore

import (
	"time"
	"strconv"
	"github.com/zgiles/mukluk"
	"github.com/zgiles/mukluk/helpers"
)

type StoreI interface {
	KVtoMUID(key string, value string) (string, error)
	KVtoMUIDs(key string, value string) ([]string, error)
	MUID(muid string) (mukluk.Node, error)
	MUIDs(muids []string) ([]mukluk.Node, error)
	UpdateOsStep(muid string, step string) (error)
}

type StoreDBI interface {
	KVtoMUID(key string, value string) (string, error)
	KVtoMUIDs(key string, value string) ([]string, error)
	MUID(muid string) (mukluk.Node, error)
	MUIDs(muids []string) ([]mukluk.Node, error)
	Update(muid string, key string, value string) (error)
}

type store struct {
  db StoreDBI
	validkeys []string
}

func (local store) KVtoMUID(key string, value string) (string, error) {
	_, keyerr := helpers.Contains(local.validkeys, key)
	switch {
		case keyerr != nil:
			return "", keyerr
	  case key == "muid":
			return key, nil
		default:
			return local.db.KVtoMUID(key, value)
	}
	// return local.db.KVtoMUID(key, value)
}

func (local store) KVtoMUIDs(key string, value string) ([]string, error) {
	_, keyerr := helpers.Contains(local.validkeys, key)
	switch {
		case keyerr != nil:
			return []string{}, keyerr
		case key == "muid":
			return []string{key}, nil
		default:
			return local.db.KVtoMUIDs(key, value)
	}
	// return local.db.KVtoMUIDs(key, value)
}

func (local store) MUID(muid string) (mukluk.Node, error) {
	return local.db.MUID(muid)
}

func (local store) MUIDs(muid []string) ([]mukluk.Node, error) {
	return local.db.MUIDs(muid)
}

func (local store) heartBeat(muid string) (int64, error) {
	i := heartbeatnow()
	e := local.db.Update(muid, "heartbeat", strconv.FormatInt(i, 10))
	return i, e
}

// maybe step should be an int64, but it's string from the db.. so..
func (local store) UpdateOsStep(muid string, step string) (error) {
	ue := local.db.Update(muid, "os_step", step)
	if ue != nil {
		return ue
	}
	_, he := local.heartBeat(muid)
	if he != nil {
		return he
	}
	return nil
}

func heartbeatnow() (int64) {
	return time.Now().Unix()
}

func New(db1 StoreDBI) StoreI {
	validkeys := []string{"uuid", "hostname", "ipv4address", "macaddress", "muid"}
  return &store{db1, validkeys}
}
