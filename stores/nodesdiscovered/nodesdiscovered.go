package nodesdiscovered

import (
	"time"
	"strconv"
)

type NodesDiscovered struct {
	Uuid string `json:"uuid"`
	Ipv4address string `json:"ipv4address"`
	Macaddress string `json:"macaddress"`
	Surpressed bool `json:"surpressed"`
	Enrolled bool `json:"enrolled"`
	Checkincount int64 `json:"checkincount"`
	Heartbeat int64 `json:"heartbeat"`
}

type NodesDiscoveredStore interface {
  SingleKV(field string, input string) (NodesDiscovered, error)
  MultiKV(field string, input string) ([]NodesDiscovered, error)
	Insert(nd NodesDiscovered) (NodesDiscovered, error)
	CreateAndInsert(uuid string, ipv4address string, macaddress string) (NodesDiscovered, error)
	HeartBeat(uuid string) (int64, error)
	UpdateCount(uuid string) (int64, error)
}

type NodesDiscoveredStoreDB interface {
  DbSingleKV(field string, input string) (NodesDiscovered, error)
  DbMultiKV(field string, input string) ([]NodesDiscovered, error)
	DbInsert(nd NodesDiscovered) (NodesDiscovered, error)
	DbUpdateSingleKV(uuid string, key string, value string) (error)
}

type store struct {
  db NodesDiscoveredStoreDB
}

func Create(uuid string, ipv4address string, macaddress string) (NodesDiscovered) {
	nd := NodesDiscovered{
		Uuid: uuid,
		Ipv4address: ipv4address,
		Macaddress: macaddress,
		Heartbeat: heartbeatnow(),
	}
	return nd
}

func (local store) SingleKV(field string, input string) (NodesDiscovered, error) {
	return local.db.DbSingleKV(field, input)
}

func (local store) MultiKV(field string, input string) ([]NodesDiscovered, error) {
	return local.db.DbMultiKV(field, input)
}

func (local store) Insert(nd NodesDiscovered) (NodesDiscovered, error) {
	return local.db.DbInsert(nd)
}

func (local store) CreateAndInsert(uuid string, ipv4address string, macaddress string) (NodesDiscovered, error) {
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

func NewNodesDiscoveredStore(db1 NodesDiscoveredStoreDB) NodesDiscoveredStore {
  return &store{db1}
}
