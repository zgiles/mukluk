package nodesdiscoveredmysql

import (
	"log"
	// "errors"
  _ "github.com/go-sql-driver/mysql"
	"database/sql"
  "github.com/zgiles/mukluk/stores/nodesdiscovered"
)

type nodesdiscoveredmysqldb struct {
  mysqldb *sql.DB
}

func NewNodesDiscoveredMysql(mysqldb *sql.DB) *nodesdiscoveredmysqldb {
	return &nodesdiscoveredmysqldb{mysqldb}
}

func (local nodesdiscoveredmysqldb) DbSingleKV(field string, input string) (nodesdiscovered.NodesDiscovered, error) {
	answer, err := local.queryGetDiscoveredNodeByField(field, input)
	if err != nil {
		return nodesdiscovered.NodesDiscovered{}, err
	}
	return answer, nil
}

func (local nodesdiscoveredmysqldb) DbMultiKV(field string, input string) ([]nodesdiscovered.NodesDiscovered, error) {
	return local.queryGetDiscoveredNodesByField(field, input)
}


func (local nodesdiscoveredmysqldb) DbInsert(nd nodesdiscovered.NodesDiscovered) (nodesdiscovered.NodesDiscovered, error) {
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

func (local nodesdiscoveredmysqldb) DbUpdateSingleKV(uuid string, key string, value string) (error) {
	stmt, stmterr := local.mysqldb.Prepare("UPDATE `nodes_discovered` SET `" + key + "` = ? WHERE `uuid` = ?")
	if stmterr != nil {
		return stmterr
	}
	res, err := stmt.Exec(value, uuid)
	if err != nil || res == nil {
		return stmterr
	}
	return nil
}


func (local nodesdiscoveredmysqldb) queryGetDiscoveredNodeByField(field string, input string) (nodesdiscovered.NodesDiscovered, error) { // input string, field string
	fn := func(input string) (nodesdiscovered.NodesDiscovered, error) {
		n := nodesdiscovered.NodesDiscovered{}
		err := local.mysqldb.QueryRow("select uuid, ipv4address, macaddress, surpressed, enrolled, checkincount, heartbeat from nodes_discovered where " + field + " = ? limit 1", input).Scan(&n.Uuid, &n.Ipv4address, &n.Macaddress, &n.Surpressed, &n.Enrolled, &n.Checkincount, &n.Heartbeat)
		if err != nil {
			return n, err
		}
		return n, nil
	}
	return fn(input)
}


func (local nodesdiscoveredmysqldb) queryGetDiscoveredNodesByField(field string, input string) ([]nodesdiscovered.NodesDiscovered, error) { // input string, field string
	fn := func(input string) ([]nodesdiscovered.NodesDiscovered, error) {
		nl := []nodesdiscovered.NodesDiscovered{}
		rows, err := local.mysqldb.Query("select uuid, ipv4address, macaddress, surpressed, enrolled, checkincount, heartbeat from nodes_discovered where " + field + " = ?", input)
		if err != nil {
			return nl, err
		}
		defer rows.Close()
		for rows.Next() {
			n := nodesdiscovered.NodesDiscovered{}
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
