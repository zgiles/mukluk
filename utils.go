package main

import (
  "errors"
)

func contains(slice []string, item string) (bool, error) {
    set := make(map[string]struct{}, len(slice))
    for _, s := range slice {
        set[s] = struct{}{}
    }
    _, ok := set[item]
    var err error
    if ok == false { err = errors.New("slice does not contain item") } else { err = nil }
    return ok, err
}

func errorinslice( e []error ) bool {
  for _, s := range e {
    if s != nil {
      return true
    }
  }
  return false
}
