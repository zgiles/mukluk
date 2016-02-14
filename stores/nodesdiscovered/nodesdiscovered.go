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
	Insert(nd NodesDiscovered) (NodesDiscovered, error)
}

type NodesDiscoveredStoreDB interface {
  DbSingleKV(field string, input string) (NodesDiscovered, error)
  DbMultiKV(field string, input string) ([]NodesDiscovered, error)
	DbInsert(nd NodesDiscovered) (NodesDiscovered, error)
}

type store struct {
  db NodesDiscoveredStoreDB
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


func NewNodesDiscoveredStore(db1 NodesDiscoveredStoreDB) NodesDiscoveredStore {
  return &store{db1}
}
