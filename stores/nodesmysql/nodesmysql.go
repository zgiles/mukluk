package nodesmysql

import (
  "errors"
  _ "github.com/go-sql-driver/mysql"
	"database/sql"
  "github.com/zgiles/mukluk"
)

// var mysqlmuiddefinition string = "CONCAT(REPLACE(uuid, '-', ''), macaddress, REPLACE(ipv4address, '.', ''))"

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

/*
func (local nodesmysqldb) DbSingleKV(field string, input string) (mukluk.Node, error) {
	answer, err := local.queryGetNodeByField(field, input)
	if err != nil {
		return mukluk.Node{}, err
	}
	return answer, nil
}

func (local nodesmysqldb) DbMultiKV(field string, input string) ([]mukluk.Node, error) {
	return local.queryGetNodesByField(field, input)
}
*/


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

/*
// needs more work. just copied from nodesdiscoveredmysql
func (local nodesmysqldb) Insert(nd mukluk.Node) (mukluk.Node, error) {
	stmt, stmterr := local.mysqldb.Prepare("insert into `nodes` (`uuid`, `ipv4address`, `macaddress`, `heartbeat`) VALUES (?, ?, ?, ?)")
	if stmterr != nil {
		return nd, stmterr
	}
	res, err := stmt.Exec(&nd.Uuid, &nd.Ipv4address, &nd.Macaddress, &nd.Heartbeat)
	if err != nil || res == nil {
		return nd, err
	}
	return nd, nil
}
*/


/*
// DB QUERIES

// queryGetNodeByField
// Goal: Given a field and a string, return a single node object from the database that matches the field name
// Assumptions:
// * Assume that "field" input is a valid column and is unique. No checking in this function
// * Columns are correct and consistent with the type Node (above)
// * Only one row will match, and only one Node is returned. This is not checked.
// Issues:
// * "field" input is not checked. SQL Injection is possible. Previous function should check it
// * If the row isn't found an empty Node object is generated
// * Probably can get rid of the limit 1
// How: This function generates a function with the embedded field variable in the SQL command. Then executes that function
func (local nodesmysqldb) queryGetNodeByField(field string, input string) (mukluk.Node, error) { // input string, field string
	fn := func(input string) (mukluk.Node, error) {
		n := mukluk.Node{}
		err := local.mysqldb.QueryRow("select uuid, hostname, ipv4address, macaddress, os_name, os_step, node_type, oob_type, heartbeat from nodes where " + field + " = ? limit 1", input).Scan(&n.Uuid, &n.Hostname, &n.Ipv4address, &n.Macaddress, &n.Os_name, &n.Os_step, &n.Node_type, &n.Oob_type, &n.Heartbeat)
		if err != nil {
      return n, err
		}
		return n, nil
	}
	return fn(input)
}

// queryGetNodesByField
func (local nodesmysqldb) queryGetNodesByField(field string, input string) ([]mukluk.Node, error) { // input string, field string
	fn := func(input string) ([]mukluk.Node, error) {
		nl := []mukluk.Node{}
		rows, err := local.mysqldb.Query("select uuid, hostname, ipv4address, macaddress, os_name, os_step, node_type, oob_type, heartbeat from nodes where " + field + " = ?", input)
		if err != nil {
      return nl, err
		}
		defer rows.Close()
		for rows.Next() {
			n := mukluk.Node{}
			err = rows.Scan(&n.Uuid, &n.Hostname, &n.Ipv4address, &n.Macaddress, &n.Os_name, &n.Os_step, &n.Node_type, &n.Oob_type, &n.Heartbeat)
			nl = append(nl, n)
		}
		if rows.Err() != nil {
      return nl, err
		}
    return nl, nil
	}
	return fn(input)
}
*/
