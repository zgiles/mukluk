package mukluk

import (
	"encoding/json"
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

func (nd NodesDiscovered) MUID() string {
	return MUID(nd.Uuid, nd.Macaddress, nd.Ipv4address)
}

func (nd NodesDiscovered) MarshalJSON() ([]byte, error) {
	type Alias NodesDiscovered
	return json.Marshal(&struct {
		MUID string `json:"muid"`
		Alias
	}{
		MUID:     nd.MUID(),
		Alias:    (Alias)(nd),
	})
}
