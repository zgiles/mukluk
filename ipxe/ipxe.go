package ipxe

import (
  "strings"
)

func CleanUUID(i string) (string) {
  return strings.ToLower(i)
}

func CleanMac(i string) (string) {
  return strings.ToUpper(i)
}

func CleanHexHyp(i string) (string) {
  a := strings.Split(i, "-")
  b := strings.Join(a, "")
  return CleanMac(b)
}
