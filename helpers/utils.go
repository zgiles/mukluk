package helpers

import (
	"errors"
	"reflect"
	"strconv"
)

func Contains(slice []string, item string) (bool, error) {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	var err error
	if ok == false {
		err = errors.New("slice does not contain item")
	} else {
		err = nil
	}
	return ok, err
}

func Errorinslice(e []error) bool {
	for _, s := range e {
		if s != nil {
			return true
		}
	}
	return false
}

// Oh my. All to avoid having to marshall and unmarshall into a map.
// Read the source for json/encode.go and seemes safe..
// functional alternative
//
// js, marshallerr := json.Marshal(o)
// if marshallerr != nil {
//   ac.errorresponse(w, http.StatusInternalServerError)
//  return
// }
// var data map[string]string
// _ = json.Unmarshal(js, &data)
// s := data[field]
//
func ReflectStructByJSONName(o interface{}, field string) (string, error) {
	b := reflect.ValueOf(o)
	c := b.Type()
	for i := 0; i < b.NumField(); i++ {
		f := c.Field(i)
		if f.Tag.Get("json") == field {
			n := f.Name
			// return b.FieldByName(n).String(), nil
			if f.Type.String() == "string" {
				return b.FieldByName(n).String(), nil
			}
			if f.Type.String() == "int" {
				return strconv.FormatInt(b.FieldByName(n).Int(), 10), nil
			}
			if f.Type.String() == "int64" {
				return strconv.FormatInt(b.FieldByName(n).Int(), 10), nil
			}
			if f.Type.String() == "bool" {
				if b.FieldByName(n).Bool() == true {
					return "true", nil
				} else {
					return "false", nil
				}
			}
			return "", errors.New("field found but not the right type")
		}
	}
	return "", errors.New("no json field by that name")
}
