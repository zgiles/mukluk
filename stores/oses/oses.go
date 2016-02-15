package oses

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


type OsStore interface {
  SingleNameStep(os_name string, os_step string) (Os, error)
}

type OsStoreDB interface {
  DbSingleNameStep(os_name string, os_step string) (Os, error)
}

type store struct {
  db OsStoreDB
}

func (local store) SingleNameStep(field string, input string) (Os, error) {
	return local.db.DbSingleNameStep(field, input)
}

func NewOsStore(db1 OsStoreDB) OsStore {
  return &store{db1}
}
