package nodesdiscoveredstore

import (
	"time"
	"strconv"
	"github.com/zgiles/mukluk"
)

type StoreI interface {
	MUID(muid string) (mukluk.NodesDiscovered, error)
	MUIDs(muids []string) ([]mukluk.NodesDiscovered, error)
	KVtoMUID(key string, value string) (string, error)
	KVtoMUIDs(key string, value string) ([]string, error)
  // SingleKV(field string, input string) (mukluk.NodesDiscovered, error)
  // MultiKV(field string, input string) ([]mukluk.NodesDiscovered, error)
	Insert(nd mukluk.NodesDiscovered) (mukluk.NodesDiscovered, error)
	CreateAndInsert(uuid string, ipv4address string, macaddress string) (mukluk.NodesDiscovered, error)
	UpdateCount(muid string) (int64, error)
}

/*
type StoreDBI interface {
  DbSingleKV(field string, input string) (mukluk.NodesDiscovered, error)
  DbMultiKV(field string, input string) ([]mukluk.NodesDiscovered, error)
	MUID(muid string) (mukluk.NodesDiscovered, error)
	DbInsert(nd mukluk.NodesDiscovered) (mukluk.NodesDiscovered, error)
	DbUpdateSingleKV(muid string, key string, value string) (error)
}
*/

type StoreDBI interface {
	MUID(muid string) (mukluk.NodesDiscovered, error)
	MUIDs(muids []string) ([]mukluk.NodesDiscovered, error)
	KVtoMUID(key string, value string) (string, error)
	KVtoMUIDs(key string, value string) ([]string, error)
	Insert(nd mukluk.NodesDiscovered) (mukluk.NodesDiscovered, error)
	Update(muid string, key string, value string) (error)
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

/*
func (local store) SingleKV(field string, input string) (mukluk.NodesDiscovered, error) {
	muid, _ := local.db.KVtoMUID(field, input)
	return local.db.MUID(muid)
}

func (local store) MultiKV(field string, input string) ([]mukluk.NodesDiscovered, error) {
	muids, _ := local.db.KVtoMUIDs(field, input)
	return local.db.MUIDs(muids)
}
*/

func (local store) KVtoMUID(key string, value string) (string, error) {
	return local.db.KVtoMUID(key, value)
}

func (local store) KVtoMUIDs(key string, value string) ([]string, error) {
	return local.db.KVtoMUIDs(key, value)
}

func (local store) MUID(muid string) (mukluk.NodesDiscovered, error) {
	return local.db.MUID(muid)
}

func (local store) MUIDs(muid []string) ([]mukluk.NodesDiscovered, error) {
	return local.db.MUIDs(muid)
}


func (local store) Insert(nd mukluk.NodesDiscovered) (mukluk.NodesDiscovered, error) {
	return local.db.Insert(nd)
}

func (local store) CreateAndInsert(uuid string, ipv4address string, macaddress string) (mukluk.NodesDiscovered, error) {
	nd := Create(uuid, ipv4address, macaddress)
	return local.db.Insert(nd)
}

func (local store) heartBeat(muid string) (int64, error) {
	i := heartbeatnow()
	e := local.db.Update(muid, "heartbeat", strconv.FormatInt(i, 10))
	return i, e
}

func (local store) UpdateCount(muid string) (int64, error) {
	n, ne := local.db.MUID(muid)
	if ne != nil {
		return 0, ne
	}
	i := n.Checkincount + 1
	ue := local.db.Update(muid, "checkincount", strconv.FormatInt(i, 10))
	if ue != nil {
		return 0, ue
	}
	_, he := local.heartBeat(muid)
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
