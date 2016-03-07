package nodesdiscovereddb

import (
	"errors"
  _ "github.com/go-sql-driver/mysql"
	"database/sql"
  "github.com/zgiles/mukluk"
)

type nodesdiscovereddb struct {
  mysqldb *sql.DB
}

func New(mysqldb *sql.DB) *nodesdiscovereddb {
	return &nodesdiscovereddb{mysqldb}
}

func (local nodesdiscovereddb) KVtoMUID(key string, value string) (string, error) {
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

func (local nodesdiscovereddb) KVtoMUIDs(key string, value string) ([]string, error) {
	var z []string
	rows, err := local.mysqldb.Query("select " + mukluk.MUIDmysqldefinition() + " from nodes_discovered where " + key + " = ?", value)
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

func (local nodesdiscovereddb) MUID(muid string) (mukluk.NodesDiscovered, error) {
	n := mukluk.NodesDiscovered{}
	err := local.mysqldb.QueryRow("select uuid, ipv4address, macaddress, surpressed, enrolled, checkincount, heartbeat from nodes_discovered where " + mukluk.MUIDmysqldefinition() + " = ? limit 1", muid).Scan(&n.Uuid, &n.Ipv4address, &n.Macaddress, &n.Surpressed, &n.Enrolled, &n.Checkincount, &n.Heartbeat)
	if err != nil {
		return n, err
	}
	return n, nil
}

func (local nodesdiscovereddb) MUIDs(muids []string) ([]mukluk.NodesDiscovered, error) {
	nl := []mukluk.NodesDiscovered{}
	for _, muid := range muids {
		nd, nde := local.MUID(muid)
		if nde != nil {
			return nl, nde
		}
		nl = append(nl, nd)
	}
	return nl, nil
}


func (local nodesdiscovereddb) Insert(nd mukluk.NodesDiscovered) (mukluk.NodesDiscovered, error) {
	stmt, stmterr := local.mysqldb.Prepare("insert into `nodes_discovered` (`uuid`, `ipv4address`, `macaddress`, `heartbeat`) VALUES (?, ?, ?, ?)")
	if stmterr != nil {
		return nd, stmterr
	}
	res, err := stmt.Exec(&nd.Uuid, &nd.Ipv4address, &nd.Macaddress, &nd.Heartbeat)
	if err != nil || res == nil {
		return nd, err
	}
	return nd, nil
}

func (local nodesdiscovereddb) Update(muid string, key string, value string) (error) {
	stmt, stmterr := local.mysqldb.Prepare("UPDATE `nodes_discovered` SET `" + key + "` = ? WHERE " + mukluk.MUIDmysqldefinition() + " = ? LIMIT 1")
	if stmterr != nil {
		return stmterr
	}
	res, err := stmt.Exec(value, muid)
	if err != nil || res == nil {
		return stmterr
	}
	return nil
}
