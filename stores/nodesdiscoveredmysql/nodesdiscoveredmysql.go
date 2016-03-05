package nodesdiscoveredmysql

import (
	"errors"
  _ "github.com/go-sql-driver/mysql"
	"database/sql"
  "github.com/zgiles/mukluk"
)

// var mysqlmuiddefinition string = "CONCAT(REPLACE(uuid, '-', ''), macaddress, REPLACE(ipv4address, '.', ''))"

type nodesdiscoveredmysqldb struct {
  mysqldb *sql.DB
}

func New(mysqldb *sql.DB) *nodesdiscoveredmysqldb {
	return &nodesdiscoveredmysqldb{mysqldb}
}

func (local nodesdiscoveredmysqldb) KVtoMUID(key string, value string) (string, error) {
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

func (local nodesdiscoveredmysqldb) KVtoMUIDs(key string, value string) ([]string, error) {
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

/*
func (local nodesdiscoveredmysqldb) DbSingleKV(field string, input string) (mukluk.NodesDiscovered, error) {
	answer, err := local.queryGetDiscoveredNodeByField(field, input)
	if err != nil {
		return mukluk.NodesDiscovered{}, err
	}
	return answer, nil
}

func (local nodesdiscoveredmysqldb) DbMultiKV(field string, input string) ([]mukluk.NodesDiscovered, error) {
	return local.queryGetDiscoveredNodesByField(field, input)
}
*/

func (local nodesdiscoveredmysqldb) MUID(muid string) (mukluk.NodesDiscovered, error) {
	n := mukluk.NodesDiscovered{}
	err := local.mysqldb.QueryRow("select uuid, ipv4address, macaddress, surpressed, enrolled, checkincount, heartbeat from nodes_discovered where " + mukluk.MUIDmysqldefinition() + " = ? limit 1", muid).Scan(&n.Uuid, &n.Ipv4address, &n.Macaddress, &n.Surpressed, &n.Enrolled, &n.Checkincount, &n.Heartbeat)
	if err != nil {
		return n, err
	}
	return n, nil
}

func (local nodesdiscoveredmysqldb) MUIDs(muids []string) ([]mukluk.NodesDiscovered, error) {
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


func (local nodesdiscoveredmysqldb) Insert(nd mukluk.NodesDiscovered) (mukluk.NodesDiscovered, error) {
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

func (local nodesdiscoveredmysqldb) Update(muid string, key string, value string) (error) {
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

/*
func (local nodesdiscoveredmysqldb) queryGetDiscoveredNodeByField(field string, input string) (mukluk.NodesDiscovered, error) { // input string, field string
	fn := func(input string) (mukluk.NodesDiscovered, error) {
		n := mukluk.NodesDiscovered{}
		err := local.mysqldb.QueryRow("select uuid, ipv4address, macaddress, surpressed, enrolled, checkincount, heartbeat from nodes_discovered where " + field + " = ? limit 1", input).Scan(&n.Uuid, &n.Ipv4address, &n.Macaddress, &n.Surpressed, &n.Enrolled, &n.Checkincount, &n.Heartbeat)
		if err != nil {
			return n, err
		}
		return n, nil
	}
	return fn(input)
}


func (local nodesdiscoveredmysqldb) queryGetDiscoveredNodesByField(field string, input string) ([]mukluk.NodesDiscovered, error) { // input string, field string
	fn := func(input string) ([]mukluk.NodesDiscovered, error) {
		nl := []mukluk.NodesDiscovered{}
		rows, err := local.mysqldb.Query("select uuid, ipv4address, macaddress, surpressed, enrolled, checkincount, heartbeat from nodes_discovered where " + field + " = ?", input)
		if err != nil {
			return nl, err
		}
		defer rows.Close()
		for rows.Next() {
			n := mukluk.NodesDiscovered{}
			err = rows.Scan(&n.Uuid, &n.Ipv4address, &n.Macaddress, &n.Surpressed, &n.Enrolled, &n.Checkincount, &n.Heartbeat)
			nl = append(nl, n)
		}
		if rows.Err() != nil {
			log.Println(err)
      return nl, err
		}
		return nl, nil
	}
	return fn(input)
}
*/
