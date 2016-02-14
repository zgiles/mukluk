package main

import (
  "net/http"
  "encoding/json"
  "github.com/julienschmidt/httprouter"
  "github.com/gorilla/context"
)

func (ac appContext) jsonresponse(w http.ResponseWriter, js []byte, status int) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(status)
	w.Write(js)
}

func (ac appContext) textresponse(w http.ResponseWriter, s string, status int) {
  w.Header().Set("Content-Type", "text/plain")
  w.WriteHeader(status)
  w.Write([]byte(s))
}

func (ac appContext) errorresponse(w http.ResponseWriter, status int) {
  // w.Header().Set("Content-Type", "application/vnd.api+json")
  w.WriteHeader(status)
  w.Write([]byte{})
}

func (ac appContext) objectmarshaltojsonresponse(w http.ResponseWriter, o interface{}, e []error) {
	if errorinslice(e) {
		ac.errorresponse(w, http.StatusNotFound)
		return
	}
	js, marshallerr := json.Marshal(o)
	if marshallerr != nil {
		ac.errorresponse(w, http.StatusInternalServerError)
		return
	}
	ac.jsonresponse(w, js, http.StatusOK)
}

func (ac appContext) objectandfieldtotextresponse(w http.ResponseWriter, o interface{}, field string, e []error) {
  if errorinslice(e) {
		ac.errorresponse(w, http.StatusNotFound)
		return
	}
  m, merr := reflectStructByJSONName(o, field)
  if merr != nil {
    ac.errorresponse(w, http.StatusNotFound)
    return
  }
	ac.textresponse(w, m, http.StatusOK)
}


// HANDLERS

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the API URL. Please read the docs (if they exist)."))
}

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
func (ac appContext) httpGetNodeByFieldHandler(w http.ResponseWriter, r *http.Request) {
	// TODO verify inputs here
	params := context.Get(r, "params").(httprouter.Params)
	validfields := []string{ "uuid", "hostname", "ipv4address", "macaddress" }
	key := params.ByName("nodekey")
	keyvalue := params.ByName("nodekeyvalue")
  field := params.ByName("field")
	_, keyerr := contains(validfields, key)
	o, oe := ac.nodestore.SingleKV(key, keyvalue)
  if field == "" {
    ac.objectmarshaltojsonresponse(w, o, []error{ keyerr, oe } )
  } else {
    ac.objectandfieldtotextresponse(w, o, field, []error{ keyerr, oe } )
  }
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
func (ac appContext) httpGetNodesByFieldHandler(w http.ResponseWriter, r *http.Request) {
	// TODO verify inputs here
	params := context.Get(r, "params").(httprouter.Params)
	validfields := []string{ "uuid", "hostname", "ipv4address", "macaddress", "os_name", "os_step", "node_type", "oob_type" }
	key := params.ByName("nodekey")
	keyvalue := params.ByName("nodekeyvalue")
	_, keyerr := contains(validfields, key)
	o, oe := ac.nodestore.MultiKV(key, keyvalue)
	ac.objectmarshaltojsonresponse(w, o, []error{ keyerr, oe } )
}


func (ac appContext) httpGetDiscoveredNodeByFieldHandler(w http.ResponseWriter, r *http.Request) {
	// TODO verify inputs here
	params := context.Get(r, "params").(httprouter.Params)
	validfields := []string{ "uuid", "hostname", "ipv4address", "macaddress" }
	key := params.ByName("nodekey")
	keyvalue := params.ByName("nodekeyvalue")
  field := params.ByName("field")
	_, keyerr := contains(validfields, key)
	o, oe := ac.nodesdiscoveredstore.SingleKV(key, keyvalue)
  if field == "" {
    ac.objectmarshaltojsonresponse(w, o, []error{ keyerr, oe } )
  } else {
    ac.objectandfieldtotextresponse(w, o, field, []error{ keyerr, oe } )
  }
}


func (ac appContext) httpGetDiscoveredNodesByFieldHandler(w http.ResponseWriter, r *http.Request) {
	// TODO verify inputs here
	params := context.Get(r, "params").(httprouter.Params)
	validfields := []string{ "uuid", "hostname", "ipv4address", "macaddress", "surpressed", "enrolled" }
	key := params.ByName("nodekey")
	keyvalue := params.ByName("nodekeyvalue")
	_, keyerr := contains(validfields, key)
	o, oe := ac.nodesdiscoveredstore.MultiKV(key, keyvalue)
	ac.objectmarshaltojsonresponse(w, o, []error{ keyerr, oe } )
}
