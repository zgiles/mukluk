package main

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"

	"github.com/zgiles/mukluk/ipxe"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
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
	if errorinslice(e) {
		ac.errorresponse(w, http.StatusBadRequest)
		return
	}
	m, merr := reflectStructByJSONName(o, field)
	if merr != nil {
		// like marshall, internal error. missing fields shouldnt get here
		ac.errorresponse(w, http.StatusInternalServerError)
		return
	}
	ac.textresponse(w, m, http.StatusOK)
}

// HANDLERS

func errorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
	w.Write([]byte("Internal Error"))
}

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
	validfields := []string{"uuid", "hostname", "ipv4address", "macaddress"}
	key := params.ByName("nodekey")
	keyvalue := params.ByName("nodekeyvalue")
	field := params.ByName("field")
	_, keyerr := contains(validfields, key)
	o, oe := ac.nodestore.SingleKV(key, keyvalue)
	if field == "" {
		ac.objectmarshaltojsonresponse(w, o, []error{keyerr, oe})
	} else {
		ac.objectandfieldtotextresponse(w, o, field, []error{keyerr, oe})
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
	validfields := []string{"uuid", "hostname", "ipv4address", "macaddress", "os_name", "os_step", "node_type", "oob_type"}
	key := params.ByName("nodekey")
	keyvalue := params.ByName("nodekeyvalue")
	_, keyerr := contains(validfields, key)
	o, oe := ac.nodestore.MultiKV(key, keyvalue)
	ac.objectmarshaltojsonresponse(w, o, []error{keyerr, oe})
}

func (ac appContext) httpGetNodeByMyIP(w http.ResponseWriter, r *http.Request) {
	// TODO verify inputs here
	params := context.Get(r, "params").(httprouter.Params)
	field := params.ByName("field")
	ipv4address, _, iperr := net.SplitHostPort(r.RemoteAddr)
	if iperr != nil {
		ac.errorresponse(w, http.StatusBadRequest)
	}
	o, oe := ac.nodestore.SingleKV("ipv4address", ipv4address)
	if field == "" {
		ac.objectmarshaltojsonresponse(w, o, []error{iperr, oe})
	} else {
		ac.objectandfieldtotextresponse(w, o, field, []error{iperr, oe})
	}
}

func (ac appContext) httpOsNodeByMyIP(w http.ResponseWriter, r *http.Request) {
	// TODO verify inputs here
	params := context.Get(r, "params").(httprouter.Params)
	field := params.ByName("field")
	ipv4address, _, iperr := net.SplitHostPort(r.RemoteAddr)
	if iperr != nil {
		ac.errorresponse(w, http.StatusBadRequest)
	}
	n, ne := ac.nodestore.SingleKV("ipv4address", ipv4address)
	o, oe := ac.osstore.SingleNameStep(n.Os_name, strconv.FormatInt(n.Os_step, 10))
	if field == "" {
		ac.objectmarshaltojsonresponse(w, o, []error{iperr, ne, oe})
	} else {
		ac.objectandfieldtotextresponse(w, o, field, []error{iperr, ne, oe})
	}
}

func (ac appContext) httpGetDiscoveredNodeByFieldHandler(w http.ResponseWriter, r *http.Request) {
	// TODO verify inputs here
	params := context.Get(r, "params").(httprouter.Params)
	validfields := []string{"uuid", "hostname", "ipv4address", "macaddress"}
	key := params.ByName("nodekey")
	keyvalue := params.ByName("nodekeyvalue")
	field := params.ByName("field")
	_, keyerr := contains(validfields, key)
	o, oe := ac.nodesdiscoveredstore.SingleKV(key, keyvalue)
	if field == "" {
		ac.objectmarshaltojsonresponse(w, o, []error{keyerr, oe})
	} else {
		ac.objectandfieldtotextresponse(w, o, field, []error{keyerr, oe})
	}
}

func (ac appContext) httpGetDiscoveredNodesByFieldHandler(w http.ResponseWriter, r *http.Request) {
	// TODO verify inputs here
	params := context.Get(r, "params").(httprouter.Params)
	validfields := []string{"uuid", "hostname", "ipv4address", "macaddress", "surpressed", "enrolled"}
	key := params.ByName("nodekey")
	keyvalue := params.ByName("nodekeyvalue")
	_, keyerr := contains(validfields, key)
	o, oe := ac.nodesdiscoveredstore.MultiKV(key, keyvalue)
	ac.objectmarshaltojsonresponse(w, o, []error{keyerr, oe})
}

func (ac appContext) httpGetDiscoveredNodeByMyIP(w http.ResponseWriter, r *http.Request) {
	// TODO verify inputs here
	params := context.Get(r, "params").(httprouter.Params)
	field := params.ByName("field")
	ipv4address, _, iperr := net.SplitHostPort(r.RemoteAddr)
	if iperr != nil {
		ac.errorresponse(w, http.StatusBadRequest)
	}
	o, oe := ac.nodesdiscoveredstore.SingleKV("ipv4address", ipv4address)
	if field == "" {
		ac.objectmarshaltojsonresponse(w, o, []error{iperr, oe})
	} else {
		ac.objectandfieldtotextresponse(w, o, field, []error{iperr, oe})
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
	s := ipxe.UuidBoot(r.Host)
	ac.textresponse(w, s, http.StatusOK)
}

func (ac appContext) httpipxeNode(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)
	validfields := []string{"uuid", "hostname", "ipv4address", "macaddress"}
	key := params.ByName("nodekey")
	keyvalue := params.ByName("nodekeyvalue")
	_, keyerr := contains(validfields, key)
	if keyerr != nil {
		// problem with the request, reply with noop
		s := ipxe.NoopString(keyerr.Error())
		ac.textresponse(w, s, http.StatusOK)
		return
	}
	n, ne := ac.nodestore.SingleKV(key, keyvalue)
	if ne != nil {
		// node not found. Generate enrollment
		s := ipxe.Enrollmentboot(r.Host)
		ac.textresponse(w, s, http.StatusOK)
		return
	}
	o, oe := ac.osstore.SingleNameStep(n.Os_name, strconv.FormatInt(n.Os_step, 10))
	if oe != nil {
		// problem with OS, generate noop
		s := ipxe.NoopString(oe.Error())
		ac.textresponse(w, s, http.StatusOK)
		return
	}
	// should be able to generate an ipxe from the OS, it will noop if not functioning
	s := ipxe.OsBoot(o, r.Host)
	// if all successful, change to next os step
	ue := ac.nodestore.UpdateOsStep(n.Uuid, o.Next_step)
	if ue != nil {
		// problem with OS, generate noop
		s := ipxe.NoopString(ue.Error())
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

	_, luerr := ac.nodesdiscoveredstore.SingleKV("uuid", uuid)
	if luerr == nil {
		// if it is found, update the count
		_, ce := ac.nodesdiscoveredstore.UpdateCount(uuid)
		if ce != nil {
			// if updating the count didnt work, something else is wrong
			s := ipxe.NoopString(ce.Error())
			ac.textresponse(w, s, http.StatusOK)
			return
		}
	} else {
		// if not found, make it
		_, oe := ac.nodesdiscoveredstore.CreateAndInsert(uuid, ipv4address, macaddress)
		if oe != nil {
			// error making it, noop the node
			s := ipxe.NoopString(oe.Error())
			ac.textresponse(w, s, http.StatusOK)
			return
		}
	}
	// if no errors, boot locally
	s := ipxe.Localboot()
	ac.textresponse(w, s, http.StatusOK)
}
