package mukluk

type NodesDiscovered struct {
	Uuid string `json:"uuid"`
	Ipv4address string `json:"ipv4address"`
	Macaddress string `json:"macaddress"`
	Surpressed bool `json:"surpressed"`
	Enrolled bool `json:"enrolled"`
	Checkincount int64 `json:"checkincount"`
	Heartbeat int64 `json:"heartbeat"`
}
