package ipxe

import (
  "strings"
  "github.com/zgiles/mukluk"
)

func CleanUUID(i string) (string) {
  return strings.ToLower(i)
}

func CleanMac(i string) (string) {
  return strings.ToLower(i)
}

func CleanHexHyp(i string) (string) {
  a := strings.Split(i, "-")
  b := strings.Join(a, "")
  return CleanMac(b)
}

func OsBoot(o mukluk.Os, s string) (string) {
  switch o.Boot_mode {
    case "pxe":
      return pxeLinuxBoot(o, s)
    case "local":
      return Localboot()
    case "holdandwait":
      return holdandwait(s)
    default:
      return Noop()
  }
}

func IdBoot(method string, s string ) (string) {
  switch method {
    case "uuid":
      r := `#!ipxe
chain http://` + s + `/api/1/node/uuid/${uuid}/ipxe
`
      return r
    case "mac":
      r := `#!ipxe
chain http://` + s + `/api/1/node/macaddress/${mac:hexhyp}/ipxe
`
      return r
    case "muid":
      r := `#!ipxe
chain http://` + s + `/api/1/node/muid/${uuid}${mac:hexhyp}${ip}/ipxe
`
      return r
    default:
      return Noop()
  }
}

func UuidBoot(s string) (string) {
  r := `#!ipxe
chain http://` + s + `/api/1/node/uuid/${uuid}/ipxe
`
  return r
}

func pxeLinuxBoot(o mukluk.Os, s string) (string) {
  r := `#!ipxe
echo generateSimpleLinuxNetBoot running...
initrd ` + o.Boot_initrd + `
kernel ` + o.Boot_kernel + ` ` + o.Boot_options + `
boot || goto error
goto exit
:error
echo Failed. Sleeping 300. then local boot.
sleep 300
goto exit
:exit
exit
`
  return r
}

func Localboot() (string) {
  r := `#!ipxe
echo local boot.
sleep 1
goto exit
:exit
exit
`
  return r
}

func holdandwait(s string) (string) {
  r := `#!ipxe
echo holdandwait...
:sleep
echo sleeping 30
sleep 30
imgfetch -n yesno http://` + s + `/api/1/me/holdandwaityesno || goto sleep
imgfree yesno
goto exit
:exit
echo rebooting
sleep 30
reboot
`
  return r
}

func Noop() (string) {
  r := `#!ipxe
echo error on server. sleep 30 then local boot.
sleep 30
goto exit
:exit
exit
`
  return r
}

func NoopString(e string) (string) {
  r := `#!ipxe
echo ` + e +`
echo error on server. sleep 30 then local boot.
sleep 30
goto exit
:exit
exit
`
  return r
}

func Enrollmentboot(s string) (string) {
  r := `#!ipxe
echo generateEnrollmentBoot...
chain http://` + s + `/api/1/discover/uuid/${uuid}/ipv4address/${ip:ipv4}/macaddress/${mac:hexhyp} || goto error
goto exit
:error
echo Failed. Sleeping 300. then local boot.
sleep 300
goto exit
:exit
exit
`
  return r
}


func ResponseDecision(action string, s string) (string) {
  switch action {
    case "noop":
      return Noop()
    case "noopstring":
      return NoopString(s)
    case "holdandwait":
      return holdandwait(s)
    case "uuidboot":
      return UuidBoot(s)
    case "localboot":
      return Localboot()
    default:
      return NoopString("ipxe.ResponseDecision switch defaulted")
  }
}
