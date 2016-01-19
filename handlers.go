package main

import (
//	"log"
	"net/http"
	"encoding/json"
  "github.com/julienschmidt/httprouter"
)


// HANDLERS

// httpGetNodeByFieldHandler
// Goal: As an HTTP Handler, take a URL via the Request and Params and go lookup the node which matches it
// additionally, make sure the URL is matching on of the chosen unique field names
// Assumptions:
// * Assuming the HTTPRouter is correctly filling in the nodekey
// * Assuming two variables are given via the Params
// Issues:
// * "contains" function needs work
// * fields are hardcoded
// * A not found node will return an empty object. We don't error to the client
// How: The nodekey is checked against validfields. If it is not contained with in the slice, an error is returned. If it is contained,
// we call the queryGetNodeByField function with the two fields. Lastly, we marshal the data into a json output via the type Node.
func (ac appContext) httpGetNodeByFieldHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO verify inputs here
	validfields := []string{ "uuid", "hostname", "ipv4address", "macaddress" }
	key := ps.ByName("nodekey")
	keyvalue := ps.ByName("nodekeyvalue")
	keyisvalid := contains(validfields, key)
	if keyisvalid != true {
		http.Error(w, "Invalid Field", http.StatusBadRequest)
		return
	}
	/*
	n := ac.queryGetNodeByField(key, keyvalue)
	*/
	n, dberr := ac.redisgetNodesByField(key, keyvalue)
	if dberr != nil || len(n) == 0 {
		http.Error(w, "Not Found", http.StatusBadRequest)
		return
	}
	js, marshallerr := json.Marshal(n[0])
	if marshallerr != nil {
		http.Error(w, marshallerr.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.Write(js)
}

// httpGetNodesByFieldHandler
// Goal: As an HTTP Handler, take a URL via the Request and Params and go lookup the nodes which matches it
// additionally, make sure the URL is matching on of the chosen unique field names (not yet). Then return a slice of nodes
// that match that need
// Assumptions:
// * Assuming the HTTPRouter is correctly filling in the nodekey
// * Assuming two variables are given via the Params
// Issues:
// * "contains" function needs work
// * fields are hardcoded
// * A not found node will return an empty object. We don't error to the client
// How: The nodekey is checked against validfields. If it is not contained with in the slice, an error is returned. If it is contained,
// we call the queryGetNodesByField function with the two fields. Lastly, we marshal the data into a json output via the type Node.
func (ac appContext) httpGetNodesByFieldHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO verify inputs here
	// validfields := []string{ "uuid", "hostname", "ipv4address", "macaddress" }
	key := ps.ByName("nodekey")
	keyvalue := ps.ByName("nodekeyvalue")
	/*
	keyisvalid := contains(validfields, key)
	if keyisvalid != true {
		http.Error(w, "Invalid Field", http.StatusBadRequest)
		return
	}
	*/
	/*
	nl := ac.queryGetNodesByField(key, keyvalue)
	*/
	nl, dberr := ac.redisgetNodesByField(key, keyvalue)
	// should we return an empty array??
	if dberr != nil || len(nl) == 0 {
		http.Error(w, "Not Found", http.StatusBadRequest)
		return
	}
	js, marshallerr := json.Marshal(nl)
	if marshallerr != nil {
		http.Error(w, marshallerr.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.Write(js)
}


func (ac appContext) httpGetDiscoveredNodeByFieldHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO verify inputs here
	validfields := []string{ "uuid", "hostname", "ipv4address", "macaddress" }
	key := ps.ByName("nodekey")
	keyvalue := ps.ByName("nodekeyvalue")
	keyisvalid := contains(validfields, key)
	if keyisvalid != true {
		http.Error(w, "Invalid Field", http.StatusBadRequest)
		return
	}
	/*
	n := ac.queryGetDiscoveredNodeByField(key, keyvalue)
	*/
	n, dberr := ac.redisgetDiscoveredNodesByField(key, keyvalue)
	if dberr != nil || len(n) == 0 {
		http.Error(w, "Not Found", http.StatusBadRequest)
		return
	}
	js, marshallerr := json.Marshal(n[0])
	if marshallerr != nil {
		http.Error(w, marshallerr.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.Write(js)
}

func (ac appContext) httpGetDiscoveredNodesByFieldHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO verify inputs here
	// validfields := []string{ "uuid", "hostname", "ipv4address", "macaddress" }
	key := ps.ByName("nodekey")
	keyvalue := ps.ByName("nodekeyvalue")
	/*
	keyisvalid := contains(validfields, key)
	if keyisvalid != true {
		http.Error(w, "Invalid Field", http.StatusBadRequest)
		return
	}
	*/
	/*
	nl := ac.queryGetDiscoveredNodesByField(key, keyvalue)
	*/
	nl, dberr := ac.redisgetDiscoveredNodesByField(key, keyvalue)
	if dberr != nil || len(nl) == 0 {
		http.Error(w, "Not Found", http.StatusBadRequest)
		return
	}
	js, marshallerr := json.Marshal(nl)
	if marshallerr != nil {
		http.Error(w, marshallerr.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.Write(js)
}
