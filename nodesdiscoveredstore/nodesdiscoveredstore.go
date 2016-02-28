package nodesdiscoveredstore

import (
	"time"
	"strconv"
	"github.com/zgiles/mukluk"
)

type StoreI interface {
  SingleKV(field string, input string) (mukluk.NodesDiscovered, error)
  MultiKV(field string, input string) ([]mukluk.NodesDiscovered, error)
	Insert(nd mukluk.NodesDiscovered) (mukluk.NodesDiscovered, error)
	CreateAndInsert(uuid string, ipv4address string, macaddress string) (mukluk.NodesDiscovered, error)
	HeartBeat(uuid string) (int64, error)
	UpdateCount(uuid string) (int64, error)
}

type StoreDBI interface {
  DbSingleKV(field string, input string) (mukluk.NodesDiscovered, error)
  DbMultiKV(field string, input string) ([]mukluk.NodesDiscovered, error)
	DbInsert(nd mukluk.NodesDiscovered) (mukluk.NodesDiscovered, error)
	DbUpdateSingleKV(uuid string, key string, value string) (error)
}

type store struct {
  db StoreDBI
}

func Create(uuid string, ipv4address string, macaddress string) (mukluk.NodesDiscovered) {
	nd := mukluk.NodesDiscovered{
		Uuid: uuid,
		Ipv4address: ipv4address,
		Macaddress: macaddress,
		Heartbeat: heartbeatnow(),
	}
	return nd
}

func (local store) SingleKV(field string, input string) (mukluk.NodesDiscovered, error) {
	return local.db.DbSingleKV(field, input)
}

func (local store) MultiKV(field string, input string) ([]mukluk.NodesDiscovered, error) {
	return local.db.DbMultiKV(field, input)
}

func (local store) Insert(nd mukluk.NodesDiscovered) (mukluk.NodesDiscovered, error) {
	return local.db.DbInsert(nd)
}

func (local store) CreateAndInsert(uuid string, ipv4address string, macaddress string) (mukluk.NodesDiscovered, error) {
	nd := Create(uuid, ipv4address, macaddress)
	return local.db.DbInsert(nd)
}

func (local store) HeartBeat(uuid string) (int64, error) {
	i := heartbeatnow()
	e := local.db.DbUpdateSingleKV(uuid, "heartbeat", strconv.FormatInt(i, 10))
	return i, e
}

func (local store) UpdateCount(uuid string) (int64, error) {
	n, ne := local.db.DbSingleKV("uuid", uuid)
	if ne != nil {
		return 0, ne
	}
	i := n.Checkincount + 1
	ue := local.db.DbUpdateSingleKV(uuid, "checkincount", strconv.FormatInt(i, 10))
	if ue != nil {
		return 0, ue
	}
	_, he := local.HeartBeat(uuid)
	if he != nil {
		return 0, he
	}
	return i, nil
}

func heartbeatnow() (int64) {
	return time.Now().Unix()
}

/*
func (local store) HeartBeatNode(uuid string) (int, error) {
	return local.db.DbHeartBeatNode(uuid)
}
*/

func New(db1 StoreDBI) StoreI {
  return &store{db1}
}
