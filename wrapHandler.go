package main

import (
  "net/http"
  "github.com/julienschmidt/httprouter"
  "github.com/gorilla/context"
)

func wrapHandler(h http.Handler) httprouter.Handle {
  return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    context.Set(r, "params", ps)
    h.ServeHTTP(w, r)
  }
}
