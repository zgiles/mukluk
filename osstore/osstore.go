package osstore

import (
	"github.com/zgiles/mukluk"
)

type StoreI interface {
  SingleNameStep(os_name string, os_step string) (mukluk.Os, error)
}

type StoreDBI interface {
  DbSingleNameStep(os_name string, os_step string) (mukluk.Os, error)
}

type store struct {
  db StoreDBI
}

func (local store) SingleNameStep(field string, input string) (mukluk.Os, error) {
	return local.db.DbSingleNameStep(field, input)
}

func New(db1 StoreDBI) StoreI {
  return &store{db1}
}
