package main

import (
	"net"
	"net/http"
	"strconv"

	"github.com/zgiles/mukluk"
	"github.com/zgiles/mukluk/ipxe"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

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
	if muiderr != nil {	errorresponse(w, http.StatusBadRequest) }
	o, oe := ac.nodestore.MUID(muid)
	switchresponsefieldornot(w, o, muid, field, []error{oe, muiderr})
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
	if muiderr != nil {	errorresponse(w, http.StatusBadRequest) }
	o, oe := ac.nodestore.MUIDs(muid)
	switchresponsefieldornot(w, o, "", "", []error{oe, muiderr})
}

func (ac appContext) httpGetNodeByMyIP(w http.ResponseWriter, r *http.Request) {
	// TODO verify inputs here
	params := context.Get(r, "params").(httprouter.Params)
	field := params.ByName("field")
	muid, muiderr := ac.mymuidbyip(r.RemoteAddr)
	if muiderr != nil {	errorresponse(w, http.StatusBadRequest) }
	o, oe := ac.nodestore.MUID(muid)
	switchresponsefieldornot(w, o, muid, field, []error{oe, muiderr})
}

func (ac appContext) httpOsNodeByMyIP(w http.ResponseWriter, r *http.Request) {
	// TODO verify inputs here
	params := context.Get(r, "params").(httprouter.Params)
	field := params.ByName("field")
	muid, muiderr := ac.mymuidbyip(r.RemoteAddr)
	if muiderr != nil {	errorresponse(w, http.StatusBadRequest) }
	n, ne := ac.nodestore.MUID(muid)
	if ne != nil {	errorresponse(w, http.StatusBadRequest) }
	o, oe := ac.osstore.SingleNameStep(n.Os_name, strconv.FormatInt(n.Os_step, 10))
	switchresponsefieldornot(w, o, "", field, []error{oe, ne, muiderr})
}

func (ac appContext) httpGetDiscoveredNodeByFieldHandler(w http.ResponseWriter, r *http.Request) {
	// TODO verify inputs here
	params := context.Get(r, "params").(httprouter.Params)
	key := params.ByName("nodekey")
	value := params.ByName("nodekeyvalue")
	field := params.ByName("field")
	if key == "macaddress" { value = ipxe.CleanHexHyp(value) }
	muid, muiderr := ac.nodesdiscoveredstore.KVtoMUID(key, value)
	if muiderr != nil {	errorresponse(w, http.StatusBadRequest) }
	o, oe := ac.nodesdiscoveredstore.MUID(muid)
	switchresponsefieldornot(w, o, muid, field, []error{oe, muiderr})
}

func (ac appContext) httpGetDiscoveredNodesByFieldHandler(w http.ResponseWriter, r *http.Request) {
	// TODO verify inputs here
	params := context.Get(r, "params").(httprouter.Params)
	key := params.ByName("nodekey")
	value := params.ByName("nodekeyvalue")
	if key == "macaddress" { value = ipxe.CleanHexHyp(value) }
	muid, muiderr := ac.nodesdiscoveredstore.KVtoMUIDs(key, value)
	if muiderr != nil {	errorresponse(w, http.StatusBadRequest) }
	o, oe := ac.nodesdiscoveredstore.MUIDs(muid)
	switchresponsefieldornot(w, o, "", "", []error{oe, muiderr})
}

func (ac appContext) httpGetDiscoveredNodeByMyIP(w http.ResponseWriter, r *http.Request) {
	// TODO verify inputs here
	params := context.Get(r, "params").(httprouter.Params)
	field := params.ByName("field")
	ipv4address, _, iperr := net.SplitHostPort(r.RemoteAddr)
	if iperr != nil {
		errorresponse(w, http.StatusBadRequest)
	}
	muid, muiderr := ac.nodesdiscoveredstore.KVtoMUID("ipv4address", ipv4address)
	if muiderr != nil {	errorresponse(w, http.StatusBadRequest) }
	o, oe := ac.nodesdiscoveredstore.MUID(muid)
	switchresponsefieldornot(w, o, "", field, []error{iperr, oe, muiderr})
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
	o, oe := ac.osstore.SingleNameStep(os_name, os_step)
	switchresponsefieldornot(w, o, "", field, []error{oe})
}

func (ac appContext) httpipxechain(w http.ResponseWriter, r *http.Request) {
	s := ipxe.IdBoot(ac.ipxeconfig.BootIDMethod, r.Host)
	textresponse(w, s, http.StatusOK)
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
		textresponse(w, s, http.StatusOK)
		return
	}
	n, ne := ac.nodestore.MUID(muid)
	if ne != nil {
		// node not found. Generate enrollment
		s := ipxe.Enrollmentboot(r.Host)
		textresponse(w, s, http.StatusOK)
		return
	}
	o, oe := ac.osstore.SingleNameStep(n.Os_name, strconv.FormatInt(n.Os_step, 10))
	if oe != nil {
		// problem with OS, do whatever we are supposed to
		s := ipxe.ResponseDecision(ac.ipxeconfig.Bootosfail, oe.Error())
		textresponse(w, s, http.StatusOK)
		return
	}
	// should be able to generate an ipxe from the OS, it will noop if not functioning
	s := ipxe.OsBoot(o, r.Host)
	// if all successful, change to next os step
	ue := ac.nodestore.UpdateOsStep(n.MUID(), o.Next_step)
	if ue != nil {
		// problem with OS, do whatever we are supposed to
		s := ipxe.ResponseDecision(ac.ipxeconfig.Bootosnextstepfail, oe.Error())
		textresponse(w, s, http.StatusOK)
		return
	}
	textresponse(w, s, http.StatusOK)
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
			textresponse(w, s, http.StatusOK)
			return
		}
		// success, go whatever we do on success
		s := ipxe.ResponseDecision(ac.ipxeconfig.Discoverandcount, "")
		textresponse(w, s, http.StatusOK)
		return
	} else {
		// if not found, make it
		_, oe := ac.nodesdiscoveredstore.CreateAndInsert(uuid, ipv4address, macaddress)
		if oe != nil {
			// error making it, do whatever we do on errors the node
			s := ipxe.ResponseDecision(ac.ipxeconfig.Discoverandinsertfail, oe.Error())
			textresponse(w, s, http.StatusOK)
			return
		}
		// success, go whatever we do on success
		s := ipxe.ResponseDecision(ac.ipxeconfig.Discoverandinsert, "")
		textresponse(w, s, http.StatusOK)
		return
	}
	// somehow we go to here, it's a bigger problem
	s := ipxe.ResponseDecision(ac.ipxeconfig.Discoverfailed, "")
	textresponse(w, s, http.StatusOK)
}