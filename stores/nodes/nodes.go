package nodes

import (
	"time"
	"strconv"
)

type Node struct {
	Uuid string `json:"uuid"`
	Hostname string `json:"hostname"`
	Ipv4address string `json:"ipv4address"`
	Macaddress string `json:"macaddress"`
	Os_name string `json:"os_name"`
	Os_step int64 `json:"os_step"`
	// Init_data []byte `json:"init_data"` // ignored because should be large
	Node_type string `json:"node_type"`
	Oob_type string `json:"oob_type"`
	Heartbeat int64 `json:"heartbeat"`
}

type NodeStore interface {
  SingleKV(field string, input string) (Node, error)
  MultiKV(field string, input string) ([]Node, error)
	HeartBeat(uuid string) (int64, error)
	UpdateOsStep(uuid string, step string) (error)
}

type NodeStoreDB interface {
  DbSingleKV(field string, input string) (Node, error)
  DbMultiKV(field string, input string) ([]Node, error)
	DbUpdateSingleKV(uuid string, key string, value string) (error)
}

type store struct {
  db NodeStoreDB
}

func (local store) SingleKV(field string, input string) (Node, error) {
	return local.db.DbSingleKV(field, input)
}

func (local store) MultiKV(field string, input string) ([]Node, error) {
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

func NewNodeStore(db1 NodeStoreDB) NodeStore {
  return &store{db1}
}
