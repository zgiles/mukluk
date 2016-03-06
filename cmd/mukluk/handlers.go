package main

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"

	"github.com/zgiles/mukluk"
	"github.com/zgiles/mukluk/helpers"
	"github.com/zgiles/mukluk/ipxe"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

var validkeys = []string{"uuid", "hostname", "ipv4address", "macaddress", "muid"}

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
	if helpers.Errorinslice(e) {
		ac.errorresponse(w, http.StatusBadRequest)
		return
	}
	js, marshallerr := json.Marshal(o)
	if marshallerr != nil {
		// marshall error is a coding or struct error, so internal server
		ac.errorresponse(w, http.StatusInternalServerError)
		return
	}
	ac.jsonresponse(w, js, http.StatusOK)
}

func (ac appContext) objectandfieldtotextresponse(w http.ResponseWriter, o interface{}, field string, e []error) {
	if helpers.Errorinslice(e) {
		ac.errorresponse(w, http.StatusBadRequest)
		return
	}
	m, merr := helpers.ReflectStructByJSONName(o, field)
	if merr != nil {
		// like marshall, internal error. missing fields shouldnt get here
		ac.errorresponse(w, http.StatusInternalServerError)
		return
	}
	ac.textresponse(w, m, http.StatusOK)
}

func (ac appContext) objecttextresponse(w http.ResponseWriter, o string, e []error) {
	if helpers.Errorinslice(e) {
		ac.errorresponse(w, http.StatusBadRequest)
		return
	}
	ac.textresponse(w, o, http.StatusOK)
}

// HANDLERS

func errorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
	w.Write([]byte("Internal Error"))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the API URL. Please read the docs (if they exist)."))
}

func (ac appContext) mymuidbyip(rawip string) (string, error) {
	ip, _, iperr := net.SplitHostPort(rawip)
	if iperr != nil { return "", iperr }
	return ac.nodestore.KVtoMUID("ipv4address", ip)
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
	key := params.ByName("nodekey")
	value := params.ByName("nodekeyvalue")
	field := params.ByName("field")
	if key == "macaddress" { value = ipxe.CleanHexHyp(value) }
	muid, muiderr := ac.nodestore.KVtoMUID(key, value)
	if muiderr != nil {	ac.errorresponse(w, http.StatusBadRequest) }
	o, oe := ac.nodestore.MUID(muid)
	switch {
	  case field == "":
			ac.objectmarshaltojsonresponse(w, o, []error{oe, muiderr})
		case field == "muid":
			ac.objecttextresponse(w, muid, []error{oe, muiderr})
		default:
			ac.objectandfieldtotextresponse(w, o, field, []error{oe, muiderr})
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
	// validfields := []string{"uuid", "hostname", "ipv4address", "macaddress", "os_name", "os_step", "node_type", "oob_type"}
	key := params.ByName("nodekey")
	value := params.ByName("nodekeyvalue")
	if key == "macaddress" { value = ipxe.CleanHexHyp(value) }
	muid, muiderr := ac.nodestore.KVtoMUIDs(key, value)
	o, oe := ac.nodestore.MUIDs(muid)
	ac.objectmarshaltojsonresponse(w, o, []error{oe, muiderr})
}

func (ac appContext) httpGetNodeByMyIP(w http.ResponseWriter, r *http.Request) {
	// TODO verify inputs here
	params := context.Get(r, "params").(httprouter.Params)
	field := params.ByName("field")
	muid, muiderr := ac.mymuidbyip(r.RemoteAddr)
	if muiderr != nil {	ac.errorresponse(w, http.StatusBadRequest) }
	o, oe := ac.nodestore.MUID(muid)
	switch true {
		case field == "":
			ac.objectmarshaltojsonresponse(w, o, []error{muiderr, oe})
		case field == "muid":
			ac.objecttextresponse(w, muid, []error{muiderr, oe})
		case muiderr == nil:
			ac.objectandfieldtotextresponse(w, o, field, []error{muiderr, oe})
		default:
			ac.errorresponse(w, http.StatusBadRequest)
	}
}

func (ac appContext) httpOsNodeByMyIP(w http.ResponseWriter, r *http.Request) {
	// TODO verify inputs here
	params := context.Get(r, "params").(httprouter.Params)
	field := params.ByName("field")
	muid, muiderr := ac.mymuidbyip(r.RemoteAddr)
	if muiderr != nil {	ac.errorresponse(w, http.StatusBadRequest) }
	n, ne := ac.nodestore.MUID(muid)
	if ne != nil {	ac.errorresponse(w, http.StatusBadRequest) }
	o, oe := ac.osstore.SingleNameStep(n.Os_name, strconv.FormatInt(n.Os_step, 10))
	if field == "" {
		ac.objectmarshaltojsonresponse(w, o, []error{muiderr, ne, oe})
	} else {
		ac.objectandfieldtotextresponse(w, o, field, []error{muiderr, ne, oe})
	}
}

func (ac appContext) httpGetDiscoveredNodeByFieldHandler(w http.ResponseWriter, r *http.Request) {
	// TODO verify inputs here
	params := context.Get(r, "params").(httprouter.Params)
	validfields := []string{"uuid", "hostname", "ipv4address", "macaddress", "muid"}
	key := params.ByName("nodekey")
	field := params.ByName("field")
	_, keyerr := helpers.Contains(validfields, key)
	var muid string
	switch key {
		case "muid":
			muid = params.ByName("nodekeyvalue")
		case "macaddress":
			keyvalue := ipxe.CleanHexHyp(params.ByName("nodekeyvalue"))
			muid, _ = ac.nodesdiscoveredstore.KVtoMUID(key, keyvalue)
		default:
			keyvalue := params.ByName("nodekeyvalue")
			muid, _ = ac.nodesdiscoveredstore.KVtoMUID(key, keyvalue)
	}
	o, oe := ac.nodesdiscoveredstore.MUID(muid)
	switch true {
	  case field == "":
			ac.objectmarshaltojsonresponse(w, o, []error{keyerr, oe})
		case field == "muid":
			ac.objecttextresponse(w, muid, []error{keyerr, oe})
		case keyerr == nil:
			ac.objectandfieldtotextresponse(w, o, field, []error{keyerr, oe})
		default:
			// field has something valid, maybe merge "contains" here and use default for error
			// ac.objectandfieldtotextresponse(w, o, field, []error{keyerr, oe})
			ac.errorresponse(w, http.StatusBadRequest)
	}
}
func (ac appContext) httpGetDiscoveredNodesByFieldHandler(w http.ResponseWriter, r *http.Request) {
	// TODO verify inputs here
	params := context.Get(r, "params").(httprouter.Params)
	key := params.ByName("nodekey")
	var muid []string
	var muiderr error
	switch key {
		case "muid":
			muid = append(muid, params.ByName("nodekeyvalue"))
		case "macaddress":
			keyvalue := ipxe.CleanHexHyp(params.ByName("nodekeyvalue"))
			muid, muiderr = ac.nodesdiscoveredstore.KVtoMUIDs(key, keyvalue)
		default:
			keyvalue := params.ByName("nodekeyvalue")
			muid, muiderr = ac.nodesdiscoveredstore.KVtoMUIDs(key, keyvalue)
	}
	o, oe := ac.nodesdiscoveredstore.MUIDs(muid)
	ac.objectmarshaltojsonresponse(w, o, []error{oe, muiderr})
}

func (ac appContext) httpGetDiscoveredNodeByMyIP(w http.ResponseWriter, r *http.Request) {
	// TODO verify inputs here
	params := context.Get(r, "params").(httprouter.Params)
	field := params.ByName("field")
	ipv4address, _, iperr := net.SplitHostPort(r.RemoteAddr)
	if iperr != nil {
		ac.errorresponse(w, http.StatusBadRequest)
	}
	muid, muiderr := ac.nodesdiscoveredstore.KVtoMUID("ipv4address", ipv4address)
	o, oe := ac.nodesdiscoveredstore.MUID(muid)
	switch true {
		case field == "":
			ac.objectmarshaltojsonresponse(w, o, []error{iperr, muiderr, oe})
		case field == "muid":
			ac.objecttextresponse(w, muid, []error{iperr, muiderr, oe})
		case oe == nil:
			ac.objectandfieldtotextresponse(w, o, field, []error{iperr, muiderr, oe})
		default:
			ac.errorresponse(w, http.StatusBadRequest)
	}
}

// httpGetOsByNameAndStepHandler
func (ac appContext) httpGetOsByNameAndStepHandler(w http.ResponseWriter, r *http.Request) {
	// TODO verify inputs here
	params := context.Get(r, "params").(httprouter.Params)
	// validfields := []string{ "uuid", "hostname", "ipv4address", "macaddress" }
	os_name := params.ByName("os_name")
	os_step := params.ByName("os_step")
	field := params.ByName("field")
	// _, keyerr := contains(validfields, key)
	var keyerr error = nil
	o, oe := ac.osstore.SingleNameStep(os_name, os_step)
	if field == "" {
		ac.objectmarshaltojsonresponse(w, o, []error{keyerr, oe})
	} else {
		ac.objectandfieldtotextresponse(w, o, field, []error{keyerr, oe})
	}
}

func (ac appContext) httpipxechain(w http.ResponseWriter, r *http.Request) {
	s := ipxe.IdBoot(ac.ipxeconfig.BootIDMethod, r.Host)
	ac.textresponse(w, s, http.StatusOK)
}

func (ac appContext) httpipxeNode(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	// validfields := []string{"uuid", "hostname", "ipv4address", "macaddress", "muid"}
	key := params.ByName("nodekey")
	value := params.ByName("nodekeyvalue")
	if key == "macaddress" { value = ipxe.CleanHexHyp(value) }
	muid, muiderr := ac.nodestore.KVtoMUID(key, value)
	if muiderr != nil {
		// problem with the request, reply with noop
		s := ipxe.ResponseDecision(ac.ipxeconfig.Badkey, muiderr.Error())
		ac.textresponse(w, s, http.StatusOK)
		return
	}
	n, ne := ac.nodestore.MUID(muid)
	if ne != nil {
		// node not found. Generate enrollment
		s := ipxe.Enrollmentboot(r.Host)
		ac.textresponse(w, s, http.StatusOK)
		return
	}
	o, oe := ac.osstore.SingleNameStep(n.Os_name, strconv.FormatInt(n.Os_step, 10))
	if oe != nil {
		// problem with OS, do whatever we are supposed to
		s := ipxe.ResponseDecision(ac.ipxeconfig.Bootosfail, oe.Error())
		ac.textresponse(w, s, http.StatusOK)
		return
	}
	// should be able to generate an ipxe from the OS, it will noop if not functioning
	s := ipxe.OsBoot(o, r.Host)
	// if all successful, change to next os step
	ue := ac.nodestore.UpdateOsStep(n.MUID(), o.Next_step)
	if ue != nil {
		// problem with OS, do whatever we are supposed to
		s := ipxe.ResponseDecision(ac.ipxeconfig.Bootosnextstepfail, oe.Error())
		ac.textresponse(w, s, http.StatusOK)
		return
	}
	ac.textresponse(w, s, http.StatusOK)
}

func (ac appContext) httpipxediscover(w http.ResponseWriter, r *http.Request) {
	// TODO verify inputs here
	params := context.Get(r, "params").(httprouter.Params)
	uuid := ipxe.CleanUUID(params.ByName("uuid"))
	ipv4address := params.ByName("ipv4address")
	macaddress := ipxe.CleanHexHyp(params.ByName("macaddress"))

	muid := mukluk.MUID(uuid, macaddress, ipv4address)

	_, luerr := ac.nodesdiscoveredstore.MUID(muid)
	if luerr == nil {
		// if it is found, update the count
		_, ce := ac.nodesdiscoveredstore.UpdateCount(muid)
		if ce != nil {
			// if updating the count didnt work, something else is wrong, return what we are supposed to
			s := ipxe.ResponseDecision(ac.ipxeconfig.Discoverandcountfail, ce.Error())
			ac.textresponse(w, s, http.StatusOK)
			return
		}
		// success, go whatever we do on success
		s := ipxe.ResponseDecision(ac.ipxeconfig.Discoverandcount, "")
		ac.textresponse(w, s, http.StatusOK)
		return
	} else {
		// if not found, make it
		_, oe := ac.nodesdiscoveredstore.CreateAndInsert(uuid, ipv4address, macaddress)
		if oe != nil {
			// error making it, do whatever we do on errors the node
			s := ipxe.ResponseDecision(ac.ipxeconfig.Discoverandinsertfail, oe.Error())
			ac.textresponse(w, s, http.StatusOK)
			return
		}
		// success, go whatever we do on success
		s := ipxe.ResponseDecision(ac.ipxeconfig.Discoverandinsert, "")
		ac.textresponse(w, s, http.StatusOK)
		return
	}
	// somehow we go to here, it's a bigger problem
	s := ipxe.ResponseDecision(ac.ipxeconfig.Discoverfailed, "")
	ac.textresponse(w, s, http.StatusOK)
}
