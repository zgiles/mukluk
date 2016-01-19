package nodesdiscovered

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
}

type NodesDiscoveredStoreDB interface {
  DbSingleKV(field string, input string) (NodesDiscovered, error)
  DbMultiKV(field string, input string) ([]NodesDiscovered, error)
}

type store struct {
  redis NodesDiscoveredStoreDB
}

func (local store) SingleKV(field string, input string) (NodesDiscovered, error) {
	return local.redis.DbSingleKV(field, input)
}

func (local store) MultiKV(field string, input string) ([]NodesDiscovered, error) {
	return local.redis.DbMultiKV(field, input)
}

func NewNodesDiscoveredStore(db1 NodesDiscoveredStoreDB) NodesDiscoveredStore {
  return &store{db1}
}
