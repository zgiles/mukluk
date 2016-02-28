package mukluk

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
