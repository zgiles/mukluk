package nodesmysql

import (
  "errors"
  _ "github.com/go-sql-driver/mysql"
	"database/sql"
  "github.com/zgiles/mukluk"
)

type nodesmysqldb struct {
  mysqldb *sql.DB
}

func New(mysqldb *sql.DB) *nodesmysqldb {
	return &nodesmysqldb{mysqldb}
}


func (local nodesmysqldb) KVtoMUID(key string, value string) (string, error) {
  a, ae := local.KVtoMUIDs(key, value)
  if ae != nil {
		return "", ae
	}
	switch len(a) {
		case 1:
			return a[0], nil
		default:
			return "", errors.New("Key Value returns more than one MUID")
	}
}

func (local nodesmysqldb) KVtoMUIDs(key string, value string) ([]string, error) {
	var z []string
	rows, err := local.mysqldb.Query("select " + mukluk.MUIDmysqldefinition() + " from nodes where " + key + " = ?", value)
	if err != nil {
		return z, err
	}
	defer rows.Close()
	for rows.Next() {
		var s string
		err = rows.Scan(&s)
		if err != nil {
			return z, err
		}
		z = append(z, s)
	}
	if rows.Err() != nil {
		return z, err
	}
	return z, nil
}

func (local nodesmysqldb) MUID(muid string) (mukluk.Node, error) {
	n := mukluk.Node{}
  err := local.mysqldb.QueryRow("select uuid, hostname, ipv4address, macaddress, os_name, os_step, node_type, oob_type, heartbeat from nodes where " + mukluk.MUIDmysqldefinition() + " = ? limit 1", muid).Scan(&n.Uuid, &n.Hostname, &n.Ipv4address, &n.Macaddress, &n.Os_name, &n.Os_step, &n.Node_type, &n.Oob_type, &n.Heartbeat)
	if err != nil {
		return n, err
	}
	return n, nil
}

func (local nodesmysqldb) MUIDs(muids []string) ([]mukluk.Node, error) {
	nl := []mukluk.Node{}
	for _, muid := range muids {
		nd, nde := local.MUID(muid)
		if nde != nil {
			return nl, nde
		}
		nl = append(nl, nd)
	}
	return nl, nil
}

func (local nodesmysqldb) Update(muid string, key string, value string) (error) {
	stmt, stmterr := local.mysqldb.Prepare("UPDATE `nodes` SET `" + key + "` = ? WHERE " + mukluk.MUIDmysqldefinition() + " = ? LIMIT 1")
	if stmterr != nil {
		return stmterr
	}
	res, err := stmt.Exec(value, muid)
	if err != nil || res == nil {
		return stmterr
	}
	return nil
}
