package mukluk

import (
  "strings"
)

type MUIDable interface {
  MUID() (string)
}

func MUID(uuid string, macaddress string, ip string) string {
	return muiduuid(uuid) + macaddress + muidip(ip)
}

func muiduuid(i string) (string) {
	a := strings.Split(i, "-")
	b := strings.Join(a, "")
	return b
}

func muidip(i string) (string) {
	a := strings.Split(i, ".")
	b := strings.Join(a, "")
	c := strings.Split(b, ":")
	d := strings.Join(c, "")
	return d
}

func MUIDmysqldefinition() (string) {
  return "CONCAT(REPLACE(uuid, '-', ''), macaddress, REPLACE(ipv4address, '.', ''))"
}
