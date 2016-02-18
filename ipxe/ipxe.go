package ipxe

import (
  "strings"

  "github.com/zgiles/mukluk/stores/oses"
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

func OsBoot(o oses.Os, s string) (string) {
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

func UuidBoot(s string) (string) {
  r := `#!ipxe
chain http://` + s + `/api/1/node/uuid/${uuid}/ipxe
`
  return r
}

func pxeLinuxBoot(o oses.Os, s string) (string) {
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
