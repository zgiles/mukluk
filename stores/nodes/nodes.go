package nodes

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
}

type NodeStoreDB interface {
  DbSingleKV(field string, input string) (Node, error)
  DbMultiKV(field string, input string) ([]Node, error)
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

func NewNodeStore(db1 NodeStoreDB) NodeStore {
  return &store{db1}
}
