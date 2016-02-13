package nodesmysql

import (
  _ "github.com/go-sql-driver/mysql"
	"database/sql"
  "gomukluk/stores/nodes"
)

type nodesmysqldb struct {
  mysqldb *sql.DB
}

func NewNodesMysql(mysqldb *sql.DB) *nodesmysqldb {
	return &nodesmysqldb{mysqldb}
}

func (local nodesmysqldb) DbSingleKV(field string, input string) (nodes.Node, error) {
	answer, err := local.queryGetNodeByField(field, input)
	if err != nil {
		return nodes.Node{}, err
	}
	return answer, nil
}

func (local nodesmysqldb) DbMultiKV(field string, input string) ([]nodes.Node, error) {
	return local.queryGetNodesByField(field, input)
}


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
func (local nodesmysqldb) queryGetNodeByField(field string, input string) (nodes.Node, error) { // input string, field string
	fn := func(input string) (nodes.Node, error) {
		n := nodes.Node{}
		err := local.mysqldb.QueryRow("select uuid, hostname, ipv4address, macaddress, os_name, os_step, node_type, oob_type, heartbeat from nodes where " + field + " = ? limit 1", input).Scan(&n.Uuid, &n.Hostname, &n.Ipv4address, &n.Macaddress, &n.Os_name, &n.Os_step, &n.Node_type, &n.Oob_type, &n.Heartbeat)
		if err != nil {
      return n, err
		}
		return n, nil
	}
	return fn(input)
}

// queryGetNodesByField
func (local nodesmysqldb) queryGetNodesByField(field string, input string) ([]nodes.Node, error) { // input string, field string
	fn := func(input string) ([]nodes.Node, error) {
		nl := []nodes.Node{}
		rows, err := local.mysqldb.Query("select uuid, hostname, ipv4address, macaddress, os_name, os_step, node_type, oob_type, heartbeat from nodes where " + field + " = ?", input)
		if err != nil {
      return nl, err
		}
		defer rows.Close()
		for rows.Next() {
			n := nodes.Node{}
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
