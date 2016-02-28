package nodestore

import (
	"time"
	"strconv"
	"github.com/zgiles/mukluk"
)

type StoreI interface {
  SingleKV(field string, input string) (mukluk.Node, error)
  MultiKV(field string, input string) ([]mukluk.Node, error)
	HeartBeat(uuid string) (int64, error)
	UpdateOsStep(uuid string, step string) (error)
}

type StoreDBI interface {
  DbSingleKV(field string, input string) (mukluk.Node, error)
  DbMultiKV(field string, input string) ([]mukluk.Node, error)
	DbUpdateSingleKV(uuid string, key string, value string) (error)
}

type store struct {
  db StoreDBI
}

func (local store) SingleKV(field string, input string) (mukluk.Node, error) {
	return local.db.DbSingleKV(field, input)
}

func (local store) MultiKV(field string, input string) ([]mukluk.Node, error) {
	return local.db.DbMultiKV(field, input)
}

func (local store) HeartBeat(uuid string) (int64, error) {
	i := heartbeatnow()
	e := local.db.DbUpdateSingleKV(uuid, "heartbeat", strconv.FormatInt(i, 10))
	return i, e
}

// maybe step should be an int64, but it's string from the db.. so..
func (local store) UpdateOsStep(uuid string, step string) (error) {
	ue := local.db.DbUpdateSingleKV(uuid, "os_step", step)
	if ue != nil {
		return ue
	}
	_, he := local.HeartBeat(uuid)
	if he != nil {
		return he
	}
	return nil
}

func heartbeatnow() (int64) {
	return time.Now().Unix()
}

func New(db1 StoreDBI) StoreI {
  return &store{db1}
}
