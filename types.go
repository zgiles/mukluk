package main


type NodesDiscovered struct {
	Uuid string `json:"uuid"`
	Ipv4address string `json:"ipv4address"`
	Macaddress string `json:"macaddress"`
	Surpressed bool `json:"surpressed"`
	Enrolled bool `json:"enrolled"`
	Checkincount int64 `json:"checkincount"`
	Heartbeat int64 `json:"heartbeat"`
}

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

type Os struct {
	Os_name string `json:"os_name"`
	Os_step int64 `json:"os_step"`
	Boot_mode string `json:"boot_mode"`
	Boot_kernel string `json:"boot_kernel"`
	Boot_initrd string `json:"boot_initrd"`
	Boot_options string `json:"boot_options"`
	Next_step string `json:"next_step"`
	// Init_data []byte `json:"init_data"` // ignored because should be large
}
