package mukluk

import (
	"encoding/json"
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

func (n Node) MUID() string {
	return MUID(n.Uuid, n.Macaddress, n.Ipv4address)
}

func (n Node) MarshalJSON() ([]byte, error) {
	type Alias Node
	return json.Marshal(&struct {
		MUID string `json:"muid"`
		Alias
	}{
		MUID:     n.MUID(),
		Alias:    (Alias)(n),
	})
}
