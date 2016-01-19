package main

import (
	"log"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
)

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
func (ac appContext) queryGetNodeByField(field string, input string) Node { // input string, field string
	fn := func(input string) Node {
		n := Node{}
		err := ac.db.QueryRow("select uuid, hostname, ipv4address, macaddress, os_name, os_step, node_type, oob_type, heartbeat from nodes where " + field + " = ? limit 1", input).Scan(&n.Uuid, &n.Hostname, &n.Ipv4address, &n.Macaddress, &n.Os_name, &n.Os_step, &n.Node_type, &n.Oob_type, &n.Heartbeat)
		if err != nil && err != sql.ErrNoRows {
			// figure out how to return a non-existant node
			log.Println(err)
		}
		return n
	}
	return fn(input)
}

// queryGetNodesByField
func (ac appContext) queryGetNodesByField(field string, input string) []Node { // input string, field string
	fn := func(input string) []Node {
		nl := []Node{}
		rows, err := ac.db.Query("select uuid, hostname, ipv4address, macaddress, os_name, os_step, node_type, oob_type, heartbeat from nodes where " + field + " = ?", input)
		if err != nil {
			log.Println(err)
			return nl // how to error
		}
		defer rows.Close()
		for rows.Next() {
			n := Node{}
			err = rows.Scan(&n.Uuid, &n.Hostname, &n.Ipv4address, &n.Macaddress, &n.Os_name, &n.Os_step, &n.Node_type, &n.Oob_type, &n.Heartbeat)
			nl = append(nl, n)
		}
		if rows.Err() != nil {
			log.Println(err)
			return nl
		}
		return nl
	}
	return fn(input)
}


func (ac appContext) queryGetDiscoveredNodeByField(field string, input string) NodesDiscovered { // input string, field string
	fn := func(input string) NodesDiscovered {
		n := NodesDiscovered{}
		err := ac.db.QueryRow("select uuid, ipv4address, macaddress, surpressed, enrolled, checkincount, heartbeat from nodes_discovered where " + field + " = ? limit 1", input).Scan(&n.Uuid, &n.Ipv4address, &n.Macaddress, &n.Surpressed, &n.Enrolled, &n.Checkincount, &n.Heartbeat)
		if err != nil && err != sql.ErrNoRows {
			// figure out how to return a non-existant node
			log.Println(err)
		}
		return n
	}
	return fn(input)
}


func (ac appContext) queryGetDiscoveredNodesByField(field string, input string) []NodesDiscovered { // input string, field string
	fn := func(input string) []NodesDiscovered {
		nl := []NodesDiscovered{}
		rows, err := ac.db.Query("select uuid, ipv4address, macaddress, surpressed, enrolled, checkincount, heartbeat from nodes_discovered where " + field + " = ?", input)
		if err != nil {
			log.Println(err)
			return nl // how to error
		}
		defer rows.Close()
		for rows.Next() {
			n := NodesDiscovered{}
			err = rows.Scan(&n.Uuid, &n.Ipv4address, &n.Macaddress, &n.Surpressed, &n.Enrolled, &n.Checkincount, &n.Heartbeat)
			nl = append(nl, n)
		}
		if rows.Err() != nil {
			log.Println(err)
			return nl
		}
		return nl
	}
	return fn(input)
}
