package nodesdiscoveredstore

import (
	"time"
	"strconv"
	"github.com/zgiles/mukluk"
	"github.com/zgiles/mukluk/helpers"
)

type StoreI interface {
	MUID(muid string) (mukluk.NodesDiscovered, error)
	MUIDs(muids []string) ([]mukluk.NodesDiscovered, error)
	KVtoMUID(key string, value string) (string, error)
	KVtoMUIDs(key string, value string) ([]string, error)
	Insert(nd mukluk.NodesDiscovered) (mukluk.NodesDiscovered, error)
	CreateAndInsert(uuid string, ipv4address string, macaddress string) (mukluk.NodesDiscovered, error)
	UpdateCount(muid string) (int64, error)
}

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
	validsinglekeys []string
	validmultikeys []string
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

func (local store) KVtoMUID(key string, value string) (string, error) {
	_, keyerr := helpers.Contains(local.validsinglekeys, key)
	switch {
		case keyerr != nil:
			return "", keyerr
		case key == "muid":
			return value, nil
		default:
			return local.db.KVtoMUID(key, value)
	}
}

func (local store) KVtoMUIDs(key string, value string) ([]string, error) {
	_, keyerr := helpers.Contains(local.validmultikeys, key)
	switch {
		case keyerr != nil:
			return []string{}, keyerr
		case key == "muid":
			return []string{value}, nil
		default:
			return local.db.KVtoMUIDs(key, value)
	}
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

func New(db1 StoreDBI) StoreI {
	validsinglekeys := []string{"uuid", "ipv4address", "macaddress", "muid"}
	validmultikeys := []string{"uuid", "ipv4address", "macaddress", "muid", "enrolled", "surpressed"}
  return &store{db: db1, validsinglekeys: validsinglekeys, validmultikeys: validmultikeys}
}
