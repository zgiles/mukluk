package osdb

import (
  _ "github.com/go-sql-driver/mysql"
	"database/sql"
  "github.com/zgiles/mukluk"
)

type osdb struct {
  mysqldb *sql.DB
}

func New(mysqldb *sql.DB) *osdb {
	return &osdb{mysqldb}
}

func (local osdb) DbSingleNameStep(os_name string, os_step string) (mukluk.Os, error) {
	answer, err := local.queryGetOsByNameAndStep(os_name, os_step)
	if err != nil {
		return mukluk.Os{}, err
	}
	return answer, nil
}

// DB QUERIES

// queryGetOsByField
func (local osdb) queryGetOsByNameAndStep(os_name string, os_step string) (mukluk.Os, error) { // input string, field string
	fn := func(os_name string, os_step string) (mukluk.Os, error) {
		n := mukluk.Os{}
    // FIX HERE
		err := local.mysqldb.QueryRow("select os_name, os_step, boot_mode, boot_kernel, boot_initrd, boot_options, next_step from os where os_name = ? and os_step = ? limit 1", os_name, os_step).Scan(&n.Os_name, &n.Os_step, &n.Boot_mode, &n.Boot_kernel, &n.Boot_initrd, &n.Boot_options, &n.Next_step)
		if err != nil {
      return n, err
		}
		return n, nil
	}
	return fn(os_name, os_step)
}
