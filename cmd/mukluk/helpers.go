package main

import (
	"encoding/json"
	"net/http"
	"github.com/zgiles/mukluk"
	"github.com/zgiles/mukluk/helpers"
)

func jsonresponse(w http.ResponseWriter, js []byte, status int) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(status)
	w.Write(js)
}

func textresponse(w http.ResponseWriter, s string, status int) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	w.Write([]byte(s))
}

func errorresponse(w http.ResponseWriter, status int) {
	// w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(status)
	w.Write([]byte{})
}

func objectmarshaltojsonresponse(w http.ResponseWriter, o interface{}, e []error) {
	if helpers.Errorinslice(e) {
		errorresponse(w, http.StatusBadRequest)
		return
	}
	js, marshallerr := json.Marshal(o)
	if marshallerr != nil {
		// marshall error is a coding or struct error, so internal server
		errorresponse(w, http.StatusInternalServerError)
		return
	}
	jsonresponse(w, js, http.StatusOK)
}

func objectandfieldtotextresponse(w http.ResponseWriter, o interface{}, field string, e []error) {
	if helpers.Errorinslice(e) {
		errorresponse(w, http.StatusBadRequest)
		return
	}
	switch field {
		case "muid":
			x, ok := o.(mukluk.MUIDable)
			if ok {
				textresponse(w, x.MUID(), http.StatusOK)
			} else {
				errorresponse(w, http.StatusBadRequest)
				return
			}
		default:
			m, merr := helpers.ReflectStructByJSONName(o, field)
			if merr != nil {
				// like marshall, internal error. missing fields shouldnt get here
				errorresponse(w, http.StatusInternalServerError)
				return
			}
			textresponse(w, m, http.StatusOK)
	}
}

func objecttextresponse(w http.ResponseWriter, o string, e []error) {
	if helpers.Errorinslice(e) {
		errorresponse(w, http.StatusBadRequest)
		return
	}
	textresponse(w, o, http.StatusOK)
}

// HANDLERS

func errorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
	w.Write([]byte("Internal Error"))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the API URL. Please read the docs (if they exist)."))
}

func switchresponsefieldornot(w http.ResponseWriter, o interface{}, field string, e []error) {
	switch {
		case field == "":
			objectmarshaltojsonresponse(w, o, e)
		default:
			objectandfieldtotextresponse(w, o, field, e)
	}
}
